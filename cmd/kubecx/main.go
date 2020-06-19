package main

import (
	"fmt"
	"log"
	"os"

	"github.com/detvdl/gokube/internal/platform/kubernetes"
	"github.com/detvdl/gokube/internal/presentation"
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

	ctx, err := presentation.SelectCtx(*c)
	err = c.SetCurrentContext(ctx.Name)
	err = c.WriteFile(kubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to write modified KubeConfig: %v\n", err)
	}
}
