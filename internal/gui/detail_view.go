package gui

import (
	"encoding/json"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type DetailView struct {
	Name     string
	View     *gocui.View
	NextView *gocui.View
	PrevView *gocui.View
}

func newDetailView() *DetailView {
	return &DetailView{
		Name: "details",
	}
}

func (gui *Gui) handleDetailsPrevView(g *gocui.Gui, v *gocui.View) error {
	if _, err := gui.g.SetCurrentView(gui.detailView.PrevView.Name()); err != nil {
		return err
	}
	return nil
}

func (v *DetailView) UpdateItems(pods []*kubernetes.Pod) error {
	return nil
}

func (v *DetailView) UpdateSelected(pod *kubernetes.Pod, dy int) error {
	v.View.Clear()
	if pod != nil {
		spec, err := json.MarshalIndent(pod.Pod.Spec, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(v.View, pod.Pod.Name)
		fmt.Fprintf(v.View, "Spec: %s\n", string(spec))
		fmt.Fprintf(v.View, "Phase: %s\n", pod.Pod.Status.Phase)
		fmt.Fprintf(v.View, "Stringified: %s\n", pod.Pod.Status.String())
	} else {
		fmt.Fprintln(v.View, "No Pod Details!")
	}
	return nil
}

func (v *DetailView) GetName() string {
	return v.Name
}
