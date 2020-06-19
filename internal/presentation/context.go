package presentation

import (
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type Context struct {
	Name      string
	IsCurrent bool
}

func GetContexts(kc kubernetes.KubeConfig) ([]Context, error) {
	var cs []Context
	for name := range kc.Config.Contexts {
		c := Context{
			Name:      name,
			IsCurrent: false,
		}
		if name == kc.Config.CurrentContext {
			c.IsCurrent = true
		}
		cs = append(cs, c)

	}
	return cs, nil
}
