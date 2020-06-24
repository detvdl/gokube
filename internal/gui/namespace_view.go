package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
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
	switch direction {
	case DIRECTION_DOWN:
		gui.state.nextNamespace()
	case DIRECTION_UP:
		gui.state.prevNamespace()
	default:
		return fmt.Errorf("Could not execute movement: %d\n", direction)
	}
	return nil
}

func (gui *Gui) handleNextNamespace(g *gocui.Gui, v *gocui.View) error {
	return gui.handleNamespaceChange(g, v, DIRECTION_DOWN)
}

func (gui *Gui) handlePrevNamespace(g *gocui.Gui, v *gocui.View) error {
	return gui.handleNamespaceChange(g, v, DIRECTION_UP)
}

func (v *NamespaceView) GetName() string {
	return v.Name
}

func (v *NamespaceView) UpdateItems(items []*kubernetes.Namespace) error {
	v.View.Clear()
	return v.render(items)
}

func (v *NamespaceView) render(items []*kubernetes.Namespace) error {
	for i, ns := range items {
		if i < len(items)-1 {
			fmt.Fprintln(v.View, ns.Name)
		} else {
			fmt.Fprint(v.View, ns.Name)
		}
	}
	return nil
}

func (v *NamespaceView) UpdateSelected(item *kubernetes.Namespace, line int) error {
	return focusPoint(v.View, 0, line)
}
