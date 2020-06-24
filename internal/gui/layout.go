package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

func (gui *Gui) layout(g *gocui.Gui) error {
	g.Highlight = true
	width, height := g.Size()
	if v, err := g.SetView("contexts", 0, 0, width/2-2, height/8-1, 0); err != nil {
		gui.contextView.View = v

		v.Highlight = true
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = "contexts"
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorWhite
		ctxs, err := GetContexts(*gui.state.kubeEnv.KubeConfig)
		if err != nil {
			return err
		}
		gui.contextView.render(gui.state.kubeEnv.CurrentContext, ctxs)
	}
	if v, err := g.SetView("namespaces", 0, height/8, width/2-2, height/3-1, 0); err != nil {
		gui.namespaceView.View = v

		v.Highlight = true
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = "namespaces"
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorWhite

		if _, err := g.SetCurrentView("namespaces"); err != nil {
			return err
		}
		gui.namespaceView.render(gui.state.namespaces.Items)
	}
	if v, err := g.SetView("pods", 0, height/3, width/2-2, height-2, 0); err != nil {
		gui.podView.View = v
		gui.podView.PrevView = gui.namespaceView.View
		gui.namespaceView.NextView = v

		v.Highlight = true
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = "pods"
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorWhite
		gui.podView.render(gui.state.pods.Items)
	}
	if v, err := g.SetView("details", width/2, 0, width-1, height-2, 0); err != nil {
		v.Highlight = false
		gui.podView.NextView = v
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = "details"
		gui.detailView.View = v
		gui.detailView.PrevView = gui.podView.View
		gui.detailView.render(gui.state.pods.GetSelected())
	}

	return nil
}

func focusPoint(v *gocui.View, cx int, cy int) error {
	lineCount := v.LinesHeight()
	if cy < 0 || cy > lineCount {
		return nil
	}
	_, height := v.Size()
	ly := height - 1
	if ly == -1 {
		ly = 0
	}
	ox, oy := v.Origin()
	// if line is above origin, move origin and set cursor to zero
	// if line is below origin + height, move origin and set cursor to max
	// otherwise set cursor to value - origin
	if ly > lineCount {
		err := v.SetCursor(cx, cy)
		if err != nil {
			return fmt.Errorf("Could not set cursor: %w\n", err)
		}
		err = v.SetOrigin(ox, 0)
		if err != nil {
			return fmt.Errorf("Could not set origin: %w\n", err)
		}
	} else if cy < oy {
		err := v.SetCursor(cx, 0)
		if err != nil {
			return fmt.Errorf("Could not set cursor: %w\n", err)
		}
		err = v.SetOrigin(ox, cy)
		if err != nil {
			return fmt.Errorf("Could not set origin: %w\n", err)
		}
	} else if cy > oy+ly {
		err := v.SetCursor(cx, ly)
		if err != nil {
			return fmt.Errorf("Could not set cursor: %w\n", err)
		}
		err = v.SetOrigin(ox, cy-ly)
		if err != nil {
			return fmt.Errorf("Could not set origin: %w\n", err)
		}
	} else {
		err := v.SetCursor(cx, cy-oy)
		if err != nil {
			return fmt.Errorf("Could not set cursor: %w\n", err)
		}
	}
	return nil
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("namespaces", gocui.KeyArrowDown, gocui.ModNone, gui.handleNextNamespace); err != nil {
		return err
	}
	if err := g.SetKeybinding("namespaces", gocui.KeyArrowUp, gocui.ModNone, gui.handlePrevNamespace); err != nil {
		return err
	}
	if err := g.SetKeybinding("namespaces", gocui.KeyArrowRight, gocui.ModNone, gui.handleNamespaceNextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("pods", gocui.KeyArrowDown, gocui.ModNone, gui.handleNextPod); err != nil {
		return err
	}
	if err := g.SetKeybinding("pods", gocui.KeyArrowUp, gocui.ModNone, gui.handlePrevPod); err != nil {
		return err
	}
	if err := g.SetKeybinding("pods", gocui.KeyArrowLeft, gocui.ModNone, gui.handlePodsPrevView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
