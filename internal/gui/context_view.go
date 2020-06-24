package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type ContextView struct {
	Name     string
	View     *gocui.View
	PrevView *gocui.View
	NextView *gocui.View
}

func newContextView() *ContextView {
	return &ContextView{
		Name: "contexts",
	}
}

func (v *ContextView) render(current *kubernetes.Context, contexts []Context) error {
	for i, ctx := range contexts {
		if i < len(contexts)-1 {
			fmt.Fprintln(v.View, ctx.Name)
		} else {
			fmt.Fprint(v.View, ctx.Name)
		}
	}
	idx, ok := indexOf(current, contexts)
	if ok {
		focusPoint(v.View, 0, idx)
	}
	return nil
}

func indexOf(c *kubernetes.Context, cs []Context) (int, bool) {
	for i, v := range cs {
		if v.Name == c.Name {
			return i, true
		}
	}
	return -1, false
}
