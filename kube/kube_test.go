package kube_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suisrc/fkssl/kube"
)

func TestCurrentNamespace(t *testing.T) {
	namespace := kube.CurrentNamespace()
	t.Logf("namespace: %s", namespace)
	assert.NotNil(t, nil)
}
