package app

import (
	"github.com/go-meson/meson/event"
	"github.com/go-meson/meson/internal/binding"
	"github.com/go-meson/meson/internal/command"
	evt "github.com/go-meson/meson/internal/event"
	"github.com/go-meson/meson/internal/object"
	"github.com/go-meson/meson/menu"
)

// App Control your application’s event lifecycle.
type App struct {
	object.Object
}

//------------------------------------------------------------------------
// Methods

func Exit(code int) {
	cmd := command.MakeCallCommand(binding.ObjApp, binding.ObjAppID, "exit", code)
	command.PostMessage(&cmd)
}

func (app *App) SetApplicationMenu(menu *menu.Menu) error {
	cmd := command.MakeCallCommand(binding.ObjApp, binding.ObjAppID, "setApplicationMenu", menu)
	_, err := command.SendMessage(&cmd)
	return err
}

//------------------------------------------------------------------------
// Callbacks

func (app *App) OnWindowCloseAll(callback event.CommonCallbackHandler) {
	const en = "window-all-closed"
	evt.AddCallback(&app.Object, en, evt.CommonCallbackItem{F: callback})
}
