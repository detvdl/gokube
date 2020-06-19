package main

import (
	"fmt"
	"log"
	"os"

	"github.com/detvdl/gokube/internal/platform/kubernetes"
	"github.com/detvdl/gokube/internal/presentation"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
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
	clientConfig, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) {
		return c.Config, nil
	})
	if err != nil {
		log.Fatalf("Failed to instantiate KubeConfig object from config path: %v", err)
	}

	clientset, err := k8s.NewForConfig(clientConfig)
	if err != nil {
		log.Fatalf("Failed to create ClientSet for configuration:\n%v\n, err: %v\n", clientConfig, err)
	}
	g, err := presentation.NewGui(clientset)
	if err != nil {
		log.Fatalf("Failed to create Gui: %v\n", err)
	}
	g.Run()
}
