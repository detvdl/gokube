package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type PodView struct {
	Pods     []*kubernetes.Pod
	Selected int
	View     *gocui.View
}

func (v *PodView) SetPods(pods []*kubernetes.Pod) error {
	v.View.Clear()
	v.Pods = pods
	for i, p := range pods {
		if i < len(pods)-1 {
			fmt.Fprintln(v.View, p.Name)
		} else {
			fmt.Fprint(v.View, p.Name)
		}
	}
	return nil
}
