package gui

import (
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
	kubeEnv       *kubernetes.Environment
	state         *guiState
	namespaceView *NamespaceView
	podView       *PodView
	detailView    *DetailView
}

type guiState struct {
	pods       *models.Pods
	namespaces []*kubernetes.Namespace
	currentPod *kubernetes.Pod
}

func (s *guiState) updatePods(pods []*kubernetes.Pod) {
	s.pods.SetItems(pods)
}

func (s *guiState) updateCurrentPod(idx int) {
	s.pods.SetCurrentItem(idx)
}

func NewGui(env *kubernetes.Environment) (*Gui, error) {
	gui := &Gui{
		kubeEnv: env,
		state: &guiState{
			pods:       models.NewPods(make([]*kubernetes.Pod, 0)),
			namespaces: make([]*kubernetes.Namespace, 0),
			currentPod: nil,
		},
		namespaceView: newNamespaceView(),
		podView:       newPodView(),
		detailView:    newDetailView(),
	}
	err := gui.state.pods.Register(gui.podView, gui.detailView)
	if err != nil {
		return nil, err
	}
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
