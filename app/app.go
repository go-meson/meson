// Package app control your application’s event lifecycle.
package app

import (
	"github.com/go-meson/meson/event"
	"github.com/go-meson/meson/internal/binding"
	"github.com/go-meson/meson/internal/command"
	evt "github.com/go-meson/meson/internal/event"
	"github.com/go-meson/meson/internal/object"
)

// App is your application’s instance.
type App struct {
	object.Object
}

//------------------------------------------------------------------------
// Methods

// Exit exit meson application with exit code.
func Exit(code int) {
	cmd := command.MakeCallCommand(binding.ObjApp, binding.ObjStaticID, "exit", code)
	if err := command.PostMessage(&cmd); err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------
// Callbacks

// OnWindowCloseAll set 'window-all-closed' event handler.
//
// 'window-all-closed' emitted when all windows have been closed.
func (app *App) OnWindowCloseAll(callback event.CommonCallbackHandler) error {
	const en = "window-all-closed"
	return evt.AddCallback(&app.Object, en, evt.CommonCallbackItem{F: callback})
}
