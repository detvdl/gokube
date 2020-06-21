package main

import (
	"fmt"
	"log"
	"os"

	"github.com/detvdl/gokube/internal/gui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not locate user's home directory!")
	}
	kubeconfigPath := fmt.Sprintf("%s/.kube/config", homeDir)
	var c kubernetes.KubeConfig
	err = c.FromFile(kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not read into json: %v\n", err)
	}

	env, err := kubernetes.NewEnvironment(&c)
	g, err := gui.NewGui(env)
	if err != nil {
		log.Fatalf("Failed to create Gui: %v\n", err)
	}
	g.Run()
}
