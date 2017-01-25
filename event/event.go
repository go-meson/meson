package event

import (
	"github.com/go-meson/meson/object"
)

type CommonCallbackHandler func(object.ObjectRef)

type CommonPreventableCallbackHandler func(object.ObjectRef) bool
