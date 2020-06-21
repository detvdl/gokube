package gui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type Gui struct {
	g             *gocui.Gui
	environment   *kubernetes.Environment
	namespaceView *NamespaceView
	podView       *PodView
}

func NewGui(env *kubernetes.Environment) (*Gui, error) {
	return &Gui{
		environment:   env,
		namespaceView: &NamespaceView{make([]*kubernetes.Namespace, 0), 0, nil},
		podView:       &PodView{make([]*kubernetes.Pod, 0), 0, nil},
	}, nil
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return err
	}
	defer g.Close()

	gui.g = g
	g.Cursor = true
	g.SetManager(gocui.ManagerFunc(gui.layout))

	if err := gui.keybindings(g); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		return err
	}
	return nil
}
