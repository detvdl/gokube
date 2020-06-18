package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/detvdl/gokube/internal/kubernetes"
	"github.com/detvdl/gokube/internal/platform/tty"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v3"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not locate user's home directory!")
	}
	path := fmt.Sprintf("%s/.kube/config", homeDir)
	var c kubernetes.KubeConfig
	err = c.FromFile(path)

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   promptui.IconSelect + " {{ .Name | cyan }}",
		Inactive: "  {{if .IsCurrent}}{{.Name | green}}{{else}}{{.Name}}{{end}}",
		Selected: promptui.IconGood + " {{ .Name | faint }}",
	}
	cs, err := c.GetContexts()
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
	configString, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("Failed to marshal edited yaml: %v\n", err)
	}
	ioutil.WriteFile(path, configString, os.ModePerm)
}
