package main

import (
	"fmt"
	"log"
	"os"

	"github.com/detvdl/gokube/cmd/kubecx/gui"
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

	ctx, err := gui.SelectCtx(c)
	if err != nil {
		log.Fatalf("Selection prompt exited: %v\n", err)
	}
	err = c.SetCurrentContext(ctx.Name)
	err = c.WriteFile(kubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to write modified KubeConfig: %v\n", err)
	}
}
