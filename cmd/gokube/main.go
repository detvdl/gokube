package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/detvdl/gokube/internal/platform/kubernetes"
	"github.com/detvdl/gokube/internal/platform/tty"
	"github.com/detvdl/gokube/internal/presentation"
	"github.com/manifoldco/promptui"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not locate user's home directory!")
	}
	kubeconfigPath := fmt.Sprintf("%s/.kube/config", homeDir)
	c, err := kubernetes.FromFile(kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not read into json: %v\n", err)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   promptui.IconSelect + " {{ .Name | cyan }} {{if .IsCurrent}}{{\"(active)\" | cyan}}{{end}}",
		Inactive: "  {{if .IsCurrent}}{{.Name | green}} {{\"(active)\" | green}}{{else}}{{.Name}}{{end}}",
		Selected: promptui.IconGood + " Switched to {{ .Name | faint }}",
	}

	cs, err := presentation.GetContexts(*c)
	if err != nil {
		log.Fatalf("Could not get contexts from kubeconfig: %v", err)
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
		fmt.Printf("Prompt failed: %q\n", err)
		return
	}
	err = c.SetCurrentContext(cs[i].Name)
	err = c.ToFile(kubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to write modified KubeConfig: %v\n", err)
	}
}
