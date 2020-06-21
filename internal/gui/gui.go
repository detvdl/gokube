package gui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type Gui struct {
	g             *gocui.Gui
	kubeEnv       *kubernetes.Environment
	state         *guiState
	namespaceView *NamespaceView
	podView       *PodView
	panelViews    []PanelView
}

type guiState struct {
	pods       []*kubernetes.Pod
	namespaces []*kubernetes.Namespace
}

func NewGui(env *kubernetes.Environment) (*Gui, error) {
	gui := &Gui{
		kubeEnv: env,
		state: &guiState{
			pods:       make([]*kubernetes.Pod, 0),
			namespaces: make([]*kubernetes.Namespace, 0),
		},
		namespaceView: newNamespaceView(),
		podView:       newPodView(),
	}
	gui.panelViews = []PanelView{gui.namespaceView, gui.podView}
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
