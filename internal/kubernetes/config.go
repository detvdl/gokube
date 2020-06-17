package kubernetes

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Context struct {
	Name      string
	IsCurrent bool
}

type KubeConfig struct {
	ApiVersion string `yaml:"apiVersion,omitempty"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters,omitempty"`
	Contexts []struct {
		Context struct {
			Cluster   string `yaml:"cluster"`
			Namespace string `yaml:"namespace"`
			User      string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts,omitempty"`
	CurrentContext string                 `yaml:"current-context,omitempty"`
	Kind           string                 `yaml:"kind"`
	Preferences    map[string]interface{} `yaml:"preferences,omitempty"`
	Users          []struct {
		Name string `yaml:"name"`
		User struct {
			Exec map[string]interface{} `yaml:"exec"`
		} `yaml:"user"`
	} `yaml:"users,omitempty"`
}

func (kc *KubeConfig) FromFile(filepath string) error {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("Error reading kubeconfig: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, kc)
	if err != nil {
		return fmt.Errorf("Error unmarshalling file: %w", err)
	}
	return nil
}

func (kc KubeConfig) GetContexts() ([]Context, error) {
	cs := make([]Context, len(kc.Contexts))
	for i, ctx := range kc.Contexts {
		cs[i] = Context{
			Name:      ctx.Name,
			IsCurrent: false,
		}
		if ctx.Name == kc.CurrentContext {
			cs[i].IsCurrent = true
		}
	}
	return cs, nil
}

func (kc *KubeConfig) SetCurrentContext(context string) error {
	kc.CurrentContext = context
	return nil
}
