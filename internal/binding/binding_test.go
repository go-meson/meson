package binding

import (
	"testing"
)

func TestBindingResolveFrameworkPath(t *testing.T) {
	p, err := resolveFrameworkPath()
	if err != nil {
		t.Errorf("resolveFrameworkPath fail: %#v", err)
	}
	t.Log(p)
}

func TestBindingInitBinding(t *testing.T) {
	err := LoadBinding()
	if err != nil {
		t.Errorf("initBinding fail: %#v", err)
	}
}
