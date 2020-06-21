package gui

import (
	"fmt"
	"sync"

	"github.com/awesome-gocui/gocui"
)

type PanelView interface {
	name() string
	init(state *guiState) error
	refresh(state *guiState) error
}

func (gui *Gui) getPanelView(name string) (PanelView, error) {
	for _, p := range gui.panelViews {
		if p.name() == name {
			return p, nil
		}
	}
	return nil, nil
}

func (gui *Gui) updateAllPanelViews() error {
	var names []string
	for _, view := range gui.panelViews {
		names = append(names, view.name())
	}
	return gui.updatePanelViews(names...)
}

func (gui *Gui) updatePanelViews(contexts ...string) error {
	errChan := make(chan error)
	doneChan := make(chan bool)

	var wg sync.WaitGroup

	updateFunc := func() {
		for _, name := range contexts {
			view, err := gui.getPanelView(name)
			if err != nil {
				errChan <- fmt.Errorf("Could not get view %s: %w\n", name, err)
			}
			if view != nil {
				wg.Add(1)
				func() {
					err = view.refresh(gui.state)
					if err != nil {
						errChan <- fmt.Errorf("Could not refresh view %s: %w\n", name, err)
					}
					wg.Done()
				}()
			}
		}
		wg.Wait()
		close(doneChan)
	}
	gui.g.Update(func(g *gocui.Gui) error {
		updateFunc()
		select {
		case <-doneChan:
			return nil
		case err := <-errChan:
			close(errChan)
			return fmt.Errorf("Fatal error while refreshing views: %w\n", err)
		}
	})
	return nil
}
