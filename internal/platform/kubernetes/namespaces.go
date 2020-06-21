package kubernetes

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Namespace struct {
	Name   string
	Status string
	Age    string
}

func (e *Environment) GetNamespaces() ([]*Namespace, error) {
	ns, err := e.Clients.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to list namespaces: %w\n", err)
	}
	namespaces := make([]*Namespace, len(ns.Items))
	for i, item := range ns.Items {
		namespaces[i] = &Namespace{Name: item.Name}
	}
	return namespaces, nil
}
