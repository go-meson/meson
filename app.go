package meson

// App Control your application’s event lifecycle.
type App struct {
	object
}

func newApp() *App {
	app := &App{object: newObject(objAppID, objApp)}
	addObject(objAppID, app)
	return app
}

//------------------------------------------------------------------------
// Methods

func Exit(code int) {
	cmd := makeCallCommand(objApp, objAppID, "exit", code)
	postMessage(&cmd)
}

func (app *App) SetApplicationMenu(menu *Menu) error {
	cmd := makeCallCommand(objApp, objAppID, "setApplicationMenu", menu)
	_, err := sendMessage(&cmd)
	return err
}

//------------------------------------------------------------------------
// Callbacks

func (app *App) OnWindowCloseAll(callback CommonCallbackHandler) {
	const en = "window-all-closed"
	app.addCallback(en, commonCallbackItem{f: callback})
}
