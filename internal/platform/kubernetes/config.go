package kubernetes

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type KubeConfig struct {
	Config *clientcmdapi.Config
}

func FromFile(fp string) (*KubeConfig, error) {
	clientCfg, err := clientcmd.LoadFromFile(fp)
	if err != nil {
		return nil, fmt.Errorf("Could not load from file: %w\n", err)
	}
	return &KubeConfig{Config: clientCfg}, nil
}

func (kc KubeConfig) WriteFile(fp string) error {
	return clientcmd.WriteToFile(*kc.Config, fp)
}

func (kc *KubeConfig) SetCurrentContext(context string) error {
	kc.Config.CurrentContext = context
	return nil
}
