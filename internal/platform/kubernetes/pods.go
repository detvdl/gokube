package kubernetes

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct {
	Name      string
	Namespace string
	corev1.Pod
}

func (e *Environment) GetPods(namespace string) ([]*Pod, error) {
	p, err := e.Clients.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to list pods: %w\n", err)
	}
	pods := make([]*Pod, len(p.Items))
	for i, item := range p.Items {
		pods[i] = &Pod{Name: item.Name, Namespace: namespace, Pod: item}
	}
	return pods, nil
}
