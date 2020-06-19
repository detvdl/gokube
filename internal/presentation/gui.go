package presentation

import (
	"github.com/awesome-gocui/gocui"
	k8s "k8s.io/client-go/kubernetes"
)

type Gui struct {
	g       *gocui.Gui
	clients *k8s.Clientset
}

func NewGui(clients *k8s.Clientset) (*Gui, error) {
	return &Gui{clients: clients}, nil
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		return err
	}
	defer g.Close()

	gui.g = g
	g.Cursor = true
	g.SetManager(gocui.ManagerFunc(gui.layout))
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
