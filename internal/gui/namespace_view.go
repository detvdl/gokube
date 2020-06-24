package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type NamespaceView struct {
	Name     string
	Selected int
	View     *gocui.View
	NextView *gocui.View
}

func newNamespaceView() *NamespaceView {
	return &NamespaceView{
		Name:     "namespaces",
		Selected: 0,
		View:     nil,
	}
}

func (gui *Gui) handleNamespaceNextView(g *gocui.Gui, v *gocui.View) error {
	if _, err := gui.g.SetCurrentView(gui.namespaceView.NextView.Name()); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) handleNamespaceChange(g *gocui.Gui, v *gocui.View, direction int) error {
	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	var cond bool
	switch direction {
	case DIRECTION_DOWN:
		cond = gui.namespaceView.Selected < len(gui.state.namespaces)-1
		if cond {
			gui.namespaceView.Selected += 1
			v.MoveCursor(0, 1, false)
		}
	case DIRECTION_UP:
		cond = gui.namespaceView.Selected > 0
		if cond {
			gui.namespaceView.Selected -= 1
			v.MoveCursor(0, -1, false)
		}
	default:
		return fmt.Errorf("Could not execute movement: %d\n", direction)
	}
	if cond {
		pods, err := gui.kubeEnv.GetPods(gui.state.namespaces[gui.namespaceView.Selected].Name)
		if err != nil {
			return fmt.Errorf("Failed to fetch pods: %w\n", err)
		}
		gui.state.updatePods(pods)
		gui.state.updateCurrentPod(0)
	}
	return nil
}

func (gui *Gui) handleNextNamespace(g *gocui.Gui, v *gocui.View) error {
	return gui.handleNamespaceChange(g, v, DIRECTION_DOWN)
}

func (gui *Gui) handlePrevNamespace(g *gocui.Gui, v *gocui.View) error {
	return gui.handleNamespaceChange(g, v, DIRECTION_UP)
}
