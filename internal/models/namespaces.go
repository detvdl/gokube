package models

import (
	"fmt"

	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type NamespaceObservable interface {
	Register(observers ...NamespaceObserver) error
	Deregister(observer NamespaceObserver) error
	notifyItemsChanged()
	notifyCurrentItemChanged(int)
}

type NamespaceObserver interface {
	GetName() string
	UpdateItems([]*kubernetes.Namespace) error
	UpdateSelected(*kubernetes.Namespace, int) error
}

type Namespaces struct {
	ObserverList   []NamespaceObserver
	Items          []*kubernetes.Namespace
	CurrentItemIdx int
}

func NewNamespaces(items []*kubernetes.Namespace) *Namespaces {
	return &Namespaces{
		Items:          items,
		CurrentItemIdx: 0,
	}
}

func (n *Namespaces) GetSelected() *kubernetes.Namespace {
	if len(n.Items) > 0 {
		return n.Items[n.CurrentItemIdx]
	}
	return nil
}

func (n *Namespaces) SetItems(items []*kubernetes.Namespace) error {
	n.Items = items
	n.notifyItemsChanged()
	return nil
}

func (n *Namespaces) NextItem() error {
	if n.CurrentItemIdx < len(n.Items)-1 {
		return n.SetCurrentItem(n.CurrentItemIdx + 1)
	}
	return nil
}

func (n *Namespaces) PrevItem() error {
	if n.CurrentItemIdx > 0 {
		return n.SetCurrentItem(n.CurrentItemIdx - 1)
	}
	return nil
}

func (n *Namespaces) SetCurrentItem(idx int) error {
	n.CurrentItemIdx = idx
	n.notifyCurrentItemChanged(idx)
	return nil
}

func (n *Namespaces) Register(observers ...NamespaceObserver) error {
	for _, o := range observers {
		err := n.register(o)
		if err != nil {
			return fmt.Errorf("Failed to register observer: %v\n", o)
		}
	}
	return nil
}

func (n *Namespaces) register(observer NamespaceObserver) error {
	n.ObserverList = append(n.ObserverList, observer)
	return nil
}

func (n *Namespaces) Deregister(observer NamespaceObserver) error {
	n.ObserverList = removeNamespaceObserver(n.ObserverList, observer)
	return nil
}

func (n *Namespaces) notifyItemsChanged() {
	for _, observer := range n.ObserverList {
		observer.UpdateItems(n.Items)
	}
}

func (n *Namespaces) notifyCurrentItemChanged(dy int) {
	for _, observer := range n.ObserverList {
		observer.UpdateSelected(n.GetSelected(), dy)
	}
}

func removeNamespaceObserver(observerList []NamespaceObserver, observerToRemove NamespaceObserver) []NamespaceObserver {
	observerListLength := len(observerList)
	for i, observer := range observerList {
		if observerToRemove.GetName() == observer.GetName() {
			observerList[observerListLength-1], observerList[i] = observerList[i], observerList[observerListLength-1]
			return observerList[:observerListLength-1]
		}
	}
	return observerList
}
