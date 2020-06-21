package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type NamespaceView struct {
	Namespaces []*kubernetes.Namespace
	Selected   int
	View       *gocui.View
}

func (gui *Gui) handleNextNamespace(g *gocui.Gui, v *gocui.View) error {
	cursorDown(g, v)
	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	if gui.namespaceView.Selected < len(gui.namespaceView.Namespaces)-1 {
		gui.namespaceView.Selected += 1
		pods, err := gui.environment.GetPods(gui.namespaceView.Namespaces[gui.namespaceView.Selected].Name)
		if err != nil {
			return fmt.Errorf("Failed to fetch pods: %w\n", err)
		}
		gui.podView.SetPods(pods)
	}
	return nil
}

func (gui *Gui) handlePrevNamespace(g *gocui.Gui, v *gocui.View) error {
	cursorUp(g, v)
	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	if gui.namespaceView.Selected > 0 {
		gui.namespaceView.Selected -= 1
		pods, err := gui.environment.GetPods(gui.namespaceView.Namespaces[gui.namespaceView.Selected].Name)
		if err != nil {
			return fmt.Errorf("Failed to fetch pods: %w\n", err)
		}
		gui.podView.SetPods(pods)
	}
	return nil
}
