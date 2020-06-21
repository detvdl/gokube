package kubernetes

import (
	"fmt"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Environment struct {
	Clients          *k8s.Clientset
	KubeConfig       *KubeConfig
	CurrentContext   *Context
	CurrentNamespace *Namespace
}

func NewEnvironment(c *KubeConfig) (*Environment, error) {
	e := &Environment{KubeConfig: c}
	err := e.SetContext(c.Config.CurrentContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to set current context: %w\n", err)
	}
	return e, nil
}

func (e *Environment) SetContext(ctx string) error {
	e.CurrentContext = &Context{Name: ctx}
	e.KubeConfig.SetCurrentContext(ctx)
	clientConfig, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) {
		return e.KubeConfig.Config, nil
	})
	if err != nil {
		return fmt.Errorf("Failed to instantiate KubeConfig object from config path: %w\n", err)
	}

	clientset, err := k8s.NewForConfig(clientConfig)
	if err != nil {
		return fmt.Errorf("Failed to create ClientSet for configuration:\n%v\n, err: %w\n", clientConfig, err)
	}
	e.Clients = clientset
	e.CurrentNamespace = &Namespace{Name: "default"}
	return nil
}
