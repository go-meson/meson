package meson

import "runtime"
import "time"
import "log"
import "github.com/go-meson/meson/app"
import "github.com/go-meson/meson/internal/binding"
import "github.com/go-meson/meson/internal/object"
import "github.com/go-meson/meson/internal/command"

func init() {
	runtime.LockOSThread()
}

func newApp() *app.App {
	app := &app.App{Object: object.NewObject(binding.ObjAppID, binding.ObjApp)}
	object.AddObject(binding.ObjAppID, app)
	return app
}

func MainLoop(args []string, onInit func(*app.App)) int {
	err := binding.LoadBinding()
	if err != nil {
		log.Fatal(err)
		return -1
	}
	app := newApp()
	go func() {
		select {
		case <-binding.ReadyChannel:
			command.APIReady = true
			onInit(app)
		case <-time.After(3 * time.Second):
			log.Fatal("Waited for 3 seconds without ready signal")
		}
	}()
	return binding.RunMesonMainLoop(args)
}
