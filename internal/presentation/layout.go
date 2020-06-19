package presentation

import (
	"context"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("namespaces", 0, 0, maxX/2+7, maxY, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Clear()

		namespaces, err := gui.clients.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Failed to list pods: %v\n", err)
		}
		var output string
		for _, ns := range namespaces.Items {
			output += ns.Name + "\n"
		}
		fmt.Fprintln(v, output)
		if _, err := g.SetCurrentView("namespaces"); err != nil {
			return err
		}
	}
	return nil
}
