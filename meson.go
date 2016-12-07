package meson

import "sync/atomic"
import "errors"

import "runtime"
import "time"
import "log"

func init() {
	runtime.LockOSThread()
}

var commandID int64
var apiReady = false

func sendMessageAsync(cmd *command, handler respHandler) {
	actionID := atomic.AddInt64(&commandID, 1)
	sendMessageInternal(cmd, actionID, handler)
}

func sendMessage(cmd *command) (interface{}, error) {
	if !tryEnterSendMessage() {
		return nil, errors.New("invalid context")
	}
	defer leaveSendMessage()
	actionID := atomic.AddInt64(&commandID, 1)
	ch := getRespChan()
	sendMessageInternal(cmd, actionID, func(r *response) {
		ch <- r
	})
	resp := <-ch
	releaseRespChan(ch)
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func MainLoop(args []string, onInit func(*App)) int {
	app := newApp()
	go func() {
		select {
		case <-readyChannel:
			apiReady = true
			onInit(app)
		case <-time.After(3 * time.Second):
			log.Fatal("Waited for 3 seconds without ready signal")
		}
	}()
	return runMesonMainLoop(args)
}
