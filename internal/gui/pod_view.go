package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type PodView struct {
	Name     string
	View     *gocui.View
	NextView *gocui.View
	PrevView *gocui.View
}

func newPodView() *PodView {
	return &PodView{
		Name: "pods",
		View: nil,
	}
}

func (gui *Gui) handlePodsPrevView(g *gocui.Gui, v *gocui.View) error {
	if _, err := gui.g.SetCurrentView(gui.podView.PrevView.Name()); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) handlePodChange(g *gocui.Gui, v *gocui.View, direction int) error {
	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	switch direction {
	case DIRECTION_DOWN:
		gui.state.pods.NextItem()
	case DIRECTION_UP:
		gui.state.pods.PrevItem()
	default:
		return fmt.Errorf("Could not execute movement: %d\n", direction)
	}
	return nil
}

func (gui *Gui) handleNextPod(g *gocui.Gui, v *gocui.View) error {
	return gui.handlePodChange(g, v, DIRECTION_DOWN)
}

func (gui *Gui) handlePrevPod(g *gocui.Gui, v *gocui.View) error {
	return gui.handlePodChange(g, v, DIRECTION_UP)
}

func (v *PodView) UpdateItems(pods []*kubernetes.Pod) error {
	v.View.Clear()
	return v.render(pods)
}

func (v *PodView) render(pods []*kubernetes.Pod) error {
	for i, p := range pods {
		if i < len(pods)-1 {
			fmt.Fprintln(v.View, p.Name)
		} else {
			fmt.Fprint(v.View, p.Name)
		}
	}
	return nil
}

func (v *PodView) UpdateSelected(pod *kubernetes.Pod, line int) error {
	return focusPoint(v.View, 0, line)
}

func (v *PodView) GetName() string {
	return v.Name
}
