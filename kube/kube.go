package kube

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// CreateClient Create the kubernetes client
func CreateClient() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, errors.Wrapf(err, "error setting up cluster config")
	}

	return kubernetes.NewForConfig(config)
}
