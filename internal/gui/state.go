package gui

import (
	"fmt"

	"github.com/detvdl/gokube/internal/models"
	"github.com/detvdl/gokube/internal/platform/kubernetes"
)

type guiState struct {
	kubeEnv    *kubernetes.Environment
	pods       *models.Pods
	namespaces *models.Namespaces
	currentPod *kubernetes.Pod
}

func (s *guiState) init() error {
	namespaces, err := s.kubeEnv.GetNamespaces()
	if err != nil {
		return fmt.Errorf("Failed to list namespaces: %w\n", err)
	}
	err = s.namespaces.SetItems(namespaces)
	if err != nil {
		return fmt.Errorf("Failed to set Namespace state: %w\n", err)
	}
	s.namespaces.SetCurrentItem(0)

	pods, err := s.kubeEnv.GetPods(s.namespaces.GetSelected().Name)
	if err != nil {
		return fmt.Errorf("Failed to fetch pods: %w\n", err)
	}
	s.pods.SetItems(pods)
	s.pods.SetCurrentItem(0)

	return nil
}

func (s *guiState) refresh() error {
	return s.init()
}

func (s *guiState) nextNamespace() error {
	err := s.namespaces.NextItem()
	if err != nil {
		return err
	}
	return s.setCurrentNamespacePods()
}

func (s *guiState) prevNamespace() error {
	err := s.namespaces.PrevItem()
	if err != nil {
		return err
	}
	return s.setCurrentNamespacePods()
}

func (s *guiState) setCurrentNamespacePods() error {
	pods, err := s.kubeEnv.GetPods(s.namespaces.GetSelected().Name)
	if err != nil {
		return fmt.Errorf("Failed to fetch pods: %w\n", err)
	}
	s.pods.SetItems(pods)
	s.pods.SetCurrentItem(0)
	return nil
}
