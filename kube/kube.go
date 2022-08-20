package kube

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	ErrNotFound = errors.New("not found")
)

var (
	client kubernetes.Interface
)

func GetClient() (kubernetes.Interface, error) {
	if client != nil {
		return client, nil
	}
	var err error
	client, err = CreateClient()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// CreateClient Create the kubernetes client
func CreateClient() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, errors.Wrapf(err, "error setting up cluster config")
	}

	return kubernetes.NewForConfig(config)
}

// 获取当前命名空间
func CurrentNamespace() string {
	namespace, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	return string(namespace)
}
