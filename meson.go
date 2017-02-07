package meson

import "runtime"
import "time"
import "log"
import "github.com/go-meson/meson/internal/binding"
import "github.com/go-meson/meson/internal/command"

func init() {
	runtime.LockOSThread()
}

// MainLoop start meson application main loop. It returns an exit code to pass to App.Exit.
func MainLoop(args []string, onInit func()) int {
	err := binding.LoadBinding()
	if err != nil {
		log.Fatal(err)
		return -1
	}
	go func() {
		select {
		case <-binding.ReadyChannel:
			command.APIReady = true
			onInit()
		case <-time.After(3 * time.Second):
			log.Fatal("Waited for 3 seconds without ready signal")
		}
	}()
	return binding.RunMesonMainLoop(args)
}
