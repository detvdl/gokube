package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type PodView struct {
	Name     string
	Selected int
	View     *gocui.View
}

func newPodView() *PodView {
	return &PodView{
		Name:     "pods",
		Selected: 0,
		View:     nil,
	}
}

func (v *PodView) name() string {
	return v.Name
}

func (v *PodView) init(state *guiState) error {
	return nil
}

func (v *PodView) refresh(state *guiState) error {
	v.View.Clear()
	for i, p := range state.pods {
		if i < len(state.pods)-1 {
			fmt.Fprintln(v.View, p.Name)
		} else {
			fmt.Fprint(v.View, p.Name)
		}
	}
	return nil
}
