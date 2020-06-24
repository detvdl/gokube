package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/detvdl/gokube/internal/models"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

const (
	DIRECTION_UP = iota
	DIRECTION_DOWN
)

type Gui struct {
	g             *gocui.Gui
	state         *guiState
	contextView   *ContextView
	namespaceView *NamespaceView
	podView       *PodView
	detailView    *DetailView
}

func NewGui(env *kubernetes.Environment) (*Gui, error) {
	gui := &Gui{
		state: &guiState{
			kubeEnv:    env,
			pods:       models.NewPods(make([]*kubernetes.Pod, 0)),
			namespaces: models.NewNamespaces(make([]*kubernetes.Namespace, 0)),
			currentPod: nil,
		},
		contextView:   newContextView(),
		namespaceView: newNamespaceView(),
		podView:       newPodView(),
		detailView:    newDetailView(),
	}
	err := gui.state.init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gui state: %w\n", err)
	}
	err = gui.state.pods.Register(gui.podView, gui.detailView)
	if err != nil {
		return nil, err
	}
	err = gui.state.namespaces.Register(gui.namespaceView)
	return gui, nil
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
