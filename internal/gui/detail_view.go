package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type DetailView struct {
	Name     string
	Selected int
	View     *gocui.View
	NextView *gocui.View
	PrevView *gocui.View
}

func newDetailView() *DetailView {
	return &DetailView{
		Name:     "details",
		Selected: 0,
		View:     nil,
	}
}

func (gui *Gui) handleDetailsPrevView(g *gocui.Gui, v *gocui.View) error {
	if _, err := gui.g.SetCurrentView(gui.detailView.PrevView.Name()); err != nil {
		return err
	}
	return nil
}

func (v *DetailView) name() string {
	return v.Name
}

func (v *DetailView) init(state *guiState) error {
	return nil
}

func (v *DetailView) refresh(state *guiState) error {
	v.View.Clear()
	if state.currentPod != nil {
		spec, err := state.currentPod.Pod.Spec.Marshal()
		if err != nil {
			return err
		}
		fmt.Fprintln(v.View, state.currentPod.Name)
		fmt.Fprintf(v.View, "Spec: %s\n", string(spec))
		fmt.Fprintf(v.View, "Phase: %s\n", state.currentPod.Pod.Status.Phase)
		fmt.Fprintf(v.View, "Stringified: %s\n", state.currentPod.Pod.Status.String())
	} else {
		fmt.Fprintln(v.View, "No Pod Details!")
	}
	return nil
}
