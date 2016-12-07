package meson

import (
	"runtime"
	"sync/atomic"
)

var (
	sendLock int32
)

func lockSendMessage() {
	for {
		if atomic.CompareAndSwapInt32(&sendLock, 0, 1) {
			return
		}
		runtime.Gosched()
	}
}

func tryEnterSendMessage() bool {
	return atomic.CompareAndSwapInt32(&sendLock, 0, 1)
}

func leaveSendMessage() {
	if !atomic.CompareAndSwapInt32(&sendLock, 1, 0) {
		panic("invalid unlock timing")
	}
}
