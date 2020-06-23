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

func (gui *Gui) handleNextNamespace(g *gocui.Gui, v *gocui.View) error {
	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	if gui.namespaceView.Selected < len(gui.state.namespaces)-1 {
		gui.namespaceView.Selected += 1
		focusPoint(v, 0, gui.namespaceView.Selected)
		pods, err := gui.kubeEnv.GetPods(gui.state.namespaces[gui.namespaceView.Selected].Name)
		if err != nil {
			return fmt.Errorf("Failed to fetch pods: %w\n", err)
		}
		gui.state.pods = pods
		err = gui.updatePanelViews("pods", "details")
		if err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) handlePrevNamespace(g *gocui.Gui, v *gocui.View) error {
	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	if gui.namespaceView.Selected > 0 {
		gui.namespaceView.Selected -= 1
		focusPoint(v, 0, gui.namespaceView.Selected)
		pods, err := gui.kubeEnv.GetPods(gui.state.namespaces[gui.namespaceView.Selected].Name)
		if err != nil {
			return fmt.Errorf("Failed to fetch pods: %w\n", err)
		}
		gui.state.pods = pods
		err = gui.updatePanelViews("pods", "details")
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *NamespaceView) name() string {
	return v.Name
}

func (v *NamespaceView) render(state *guiState) error {
	return nil
}

func (v *NamespaceView) refresh(state *guiState) error {
	return nil
}
