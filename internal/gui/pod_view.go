package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type PodView struct {
	Name     string
	Selected int
	View     *gocui.View
	NextView *gocui.View
	PrevView *gocui.View
}

func newPodView() *PodView {
	return &PodView{
		Name:     "pods",
		Selected: 0,
		View:     nil,
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
	var cond bool
	switch direction {
	case DIRECTION_DOWN:
		cond = gui.podView.Selected < len(gui.state.pods)-1
		if cond {
			gui.podView.Selected += 1
			v.MoveCursor(0, 1, false)
		}
	case DIRECTION_UP:
		cond = gui.podView.Selected > 0
		if cond {
			gui.podView.Selected -= 1
			v.MoveCursor(0, -1, false)
		}
	default:
		return fmt.Errorf("Could not execute movement: %d\n", direction)
	}
	if cond {
		gui.state.currentPod = gui.state.pods[gui.podView.Selected]
		err := gui.updatePanelViews("details")
		if err != nil {
			return err
		}
	}
	return nil
}

func (gui *Gui) handleNextPod(g *gocui.Gui, v *gocui.View) error {
	return gui.handlePodChange(g, v, DIRECTION_DOWN)
}

func (gui *Gui) handlePrevPod(g *gocui.Gui, v *gocui.View) error {
	return gui.handlePodChange(g, v, DIRECTION_UP)
}

func (v *PodView) name() string {
	return v.Name
}

func (v *PodView) render(state *guiState) error {
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

func (v *PodView) refresh(state *guiState) error {
	v.View.MoveCursor(0, -(v.Selected), false)
	v.render(state)
	v.Selected = 0
	if len(state.pods) != 0 {
		state.currentPod = state.pods[v.Selected]
	} else {
		state.currentPod = nil
	}

	return nil
}
