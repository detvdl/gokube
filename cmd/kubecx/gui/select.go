package gui

import (
	"fmt"
	"strings"

	"github.com/detvdl/gokube/internal/gui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
	"github.com/detvdl/gokube/internal/platform/tty"
	"github.com/manifoldco/promptui"
)

func SelectCtx(c kubernetes.KubeConfig) (*gui.Context, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   promptui.IconSelect + " {{ .Name | cyan }} {{if .IsCurrent}}{{\"(active)\" | cyan}}{{end}}",
		Inactive: "  {{if .IsCurrent}}{{.Name | green}} {{\"(active)\" | green}}{{else}}{{.Name}}{{end}}",
		Selected: promptui.IconGood + " Switched to {{ .Name | faint }}",
	}

	cs, err := gui.GetContexts(c)
	if err != nil {
		return nil, fmt.Errorf("Could not get contexts from kubeconfig: %w", err)
	}
	searcher := func(input string, index int) bool {
		ctx := cs[index]
		name := strings.Replace(strings.ToLower(ctx.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}
	prompt := promptui.Select{
		Label:             "Select Kubernetes context",
		Items:             cs,
		Templates:         templates,
		Searcher:          searcher,
		StartInSearchMode: true,
		Size:              len(cs),
		Stdout:            &tty.BellSkipper{},
	}
	i, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w\n", err)
	}
	return &cs[i], nil
}
