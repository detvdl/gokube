package models

import (
	"fmt"

	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type PodObservable interface {
	Register(observers ...PodObserver) error
	Deregister(observer PodObserver) error
	notifyItemsChanged()
	notifyCurrentItemChanged(int)
}

type PodObserver interface {
	GetName() string
	UpdateItems([]*kubernetes.Pod) error
	UpdateSelected(*kubernetes.Pod, int) error
}

type Pods struct {
	ObserverList   []PodObserver
	Items          []*kubernetes.Pod
	CurrentItemIdx int
}

func NewPods(items []*kubernetes.Pod) *Pods {
	return &Pods{
		Items:          items,
		CurrentItemIdx: 0,
	}
}

func (p *Pods) Reset() error {
	p.SetItems(nil)
	p.SetCurrentItem(0)
	return nil
}

func (p *Pods) GetSelected() *kubernetes.Pod {
	if len(p.Items) > 0 {
		return p.Items[p.CurrentItemIdx]
	}
	return nil
}

func (p *Pods) SetItems(items []*kubernetes.Pod) error {
	p.Items = items
	p.notifyItemsChanged()
	return nil
}

func (p *Pods) NextItem() error {
	if p.CurrentItemIdx < len(p.Items)-1 {
		return p.SetCurrentItem(p.CurrentItemIdx + 1)
	}
	return nil
}

func (p *Pods) PrevItem() error {
	if p.CurrentItemIdx > 0 {
		return p.SetCurrentItem(p.CurrentItemIdx - 1)
	}
	return nil
}

func (p *Pods) SetCurrentItem(idx int) error {
	p.CurrentItemIdx = idx
	p.notifyCurrentItemChanged(idx)
	return nil
}

func (p *Pods) Register(observers ...PodObserver) error {
	for _, o := range observers {
		err := p.register(o)
		if err != nil {
			return fmt.Errorf("Failed to register observer: %v\n", o)
		}
	}
	return nil
}

func (p *Pods) register(observer PodObserver) error {
	p.ObserverList = append(p.ObserverList, observer)
	return nil
}

func (p *Pods) Deregister(observer PodObserver) error {
	p.ObserverList = removePodObserver(p.ObserverList, observer)
	return nil
}

func (p *Pods) notifyItemsChanged() {
	for _, observer := range p.ObserverList {
		observer.UpdateItems(p.Items)
	}
}

func (p *Pods) notifyCurrentItemChanged(dy int) {
	for _, observer := range p.ObserverList {
		observer.UpdateSelected(p.GetSelected(), dy)
	}
}

func removePodObserver(observerList []PodObserver, observerToRemove PodObserver) []PodObserver {
	observerListLength := len(observerList)
	for i, observer := range observerList {
		if observerToRemove.GetName() == observer.GetName() {
			observerList[observerListLength-1], observerList[i] = observerList[i], observerList[observerListLength-1]
			return observerList[:observerListLength-1]
		}
	}
	return observerList
}
