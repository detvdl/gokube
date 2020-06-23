package gui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

const (
	DIRECTION_UP = iota
	DIRECTION_DOWN
)

type Gui struct {
	g             *gocui.Gui
	kubeEnv       *kubernetes.Environment
	state         *guiState
	namespaceView *NamespaceView
	podView       *PodView
	detailView    *DetailView
	panelViews    []PanelView
}

type guiState struct {
	pods       []*kubernetes.Pod
	namespaces []*kubernetes.Namespace
	currentPod *kubernetes.Pod
}

func NewGui(env *kubernetes.Environment) (*Gui, error) {
	gui := &Gui{
		kubeEnv: env,
		state: &guiState{
			pods:       make([]*kubernetes.Pod, 0),
			namespaces: make([]*kubernetes.Namespace, 0),
			currentPod: nil,
		},
		namespaceView: newNamespaceView(),
		podView:       newPodView(),
		detailView:    newDetailView(),
	}
	gui.panelViews = []PanelView{gui.namespaceView, gui.podView, gui.detailView}
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
