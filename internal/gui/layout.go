package gui

import (
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

func (gui *Gui) layout(g *gocui.Gui) error {
	g.Highlight = true
	width, height := g.Size()
	if v, err := g.SetView("namespaces", 0, 0, width/2-2, height/3, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = "namespaces"
		v.Highlight = true
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorWhite
		namespaces, err := gui.kubeEnv.GetNamespaces()
		gui.state.namespaces = namespaces
		if err != nil {
			log.Fatalf("Failed to list pods: %v\n", err)
		}
		for i, ns := range namespaces {
			if i < len(namespaces)-1 {
				fmt.Fprintln(v, ns.Name)
			} else {
				fmt.Fprint(v, ns.Name)
			}
		}
		gui.namespaceView.View = v
		if _, err := g.SetCurrentView("namespaces"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("pods", width/2, 0, width-2, height-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = "pods"
		v.Highlight = true
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorWhite
		pods, err := gui.kubeEnv.GetPods(gui.state.namespaces[gui.namespaceView.Selected].Name)
		if err != nil {
			return fmt.Errorf("Failed to fetch pods: %w\n", err)
		}
		gui.podView.View = v
		gui.state.pods = pods
	}
	err := gui.updateAllPanelViews()
	if err != nil {
		return fmt.Errorf("Failed to update all views: %w\n", err)
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := FocusPoint(v, cx, cy+1); err != nil {
			return fmt.Errorf("Could not focus point: %w\n", err)
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := FocusPoint(v, cx, cy-1); err != nil {
			return fmt.Errorf("Could not focus point: %w\n", err)
		}
	}
	return nil
}

func FocusPoint(v *gocui.View, cx int, cy int) error {
	lineCount := len(v.BufferLines())
	if cy < 0 || cy > lineCount-1 {
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
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
