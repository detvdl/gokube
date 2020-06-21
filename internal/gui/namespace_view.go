package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type NamespaceView struct {
	Name     string
	Selected int
	View     *gocui.View
}

func newNamespaceView() *NamespaceView {
	return &NamespaceView{
		Name:     "namespaces",
		Selected: 0,
		View:     nil,
	}
}

func (gui *Gui) handleNextNamespace(g *gocui.Gui, v *gocui.View) error {
	cursorDown(g, v)
	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	if gui.namespaceView.Selected < len(gui.state.namespaces)-1 {
		gui.namespaceView.Selected += 1
		pods, err := gui.kubeEnv.GetPods(gui.state.namespaces[gui.namespaceView.Selected].Name)
		if err != nil {
			return fmt.Errorf("Failed to fetch pods: %w\n", err)
		}
		gui.state.pods = pods
	}
	gui.updatePanelViews("pods")
	return nil
}

func (gui *Gui) handlePrevNamespace(g *gocui.Gui, v *gocui.View) error {
	cursorUp(g, v)
	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	if gui.namespaceView.Selected > 0 {
		gui.namespaceView.Selected -= 1
		pods, err := gui.kubeEnv.GetPods(gui.state.namespaces[gui.namespaceView.Selected].Name)
		if err != nil {
			return fmt.Errorf("Failed to fetch pods: %w\n", err)
		}
		gui.state.pods = pods
	}
	gui.updatePanelViews("pods")
	return nil
}

func (v *NamespaceView) name() string {
	return v.Name
}

func (v *NamespaceView) init(state *guiState) error {
	return nil
}

func (v *NamespaceView) refresh(state *guiState) error {
	return nil
}