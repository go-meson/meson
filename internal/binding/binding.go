package binding

/*
#cgo darwin		CFLAGS: -mmacosx-version-min=10.9
#cgo !framework_debug	CFLAGS:	-I ./include
#cgo framework_debug	CFLAGS: -I ../../../framework/src/api
#cgo darwin LDFLAGS: -mmacosx-version-min=10.9
#cgo darwin LDFLAGS: -framework Foundation

#include <stdlib.h>
#include "binding.h"
#include "meson.h"
*/
import "C"
import "unsafe"
import "errors"
import "runtime"
import "sync/atomic"
import "github.com/go-meson/meson/object"
import "fmt"
import "reflect"

type ActionType int

const (
	ActCreate   ActionType = C.MESON_ACTION_TYPE_CREATE
	ActDelete              = C.MESON_ACTION_TYPE_DELETE
	ActCall                = C.MESON_ACTION_TYPE_CALL
	ActReply               = C.MESON_ACTION_TYPE_REPLY
	ActEvent               = C.MESON_ACTION_TYPE_EVENT
	ActRegEvent            = C.MESON_ACTION_TYPE_REGISTER_EVENT
)

const (
	ObjAppID int64 = C.MESON_OBJID_APP
)

const (
	ObjApp         object.ObjectType = C.MESON_OBJECT_TYPE_APP
	ObjWindow                        = C.MESON_OBJECT_TYPE_WINDOW
	ObjSession                       = C.MESON_OBJECT_TYPE_SESSION
	ObjWebContents                   = C.MESON_OBJECT_TYPE_WEB_CONTENTS
	ObjMenu                          = C.MESON_OBJECT_TYPE_MENU
)

type MenuType int

const (
	MenuTypeNormal    MenuType = C.MESON_MENU_TYPE_NORMAL
	MenuTypeSeparator          = C.MESON_MENU_TYPE_SEPARATOR
	MenuTypeSubmenu            = C.MESON_MENU_TYPE_SUBMENU
	MenuTypeCheckBox           = C.MESON_MENU_TYPE_CHECKBOX
	MenuTypeRadio              = C.MESON_MENU_TYPE_RADIO
)

type MessageBoxType int

const (
	MessageBoxTypeNone     MessageBoxType = C.MESON_DIALOG_MESSAGEBOX_TYPE_NONE
	MessageBoxTypeInfo                    = C.MESON_DIALOG_MESSAGEBOX_TYPE_INFO
	MessageBoxTypeWarning                 = C.MESON_DIALOG_MESSAGEBOX_TYPE_WARNING
	MessageBoxTypeError                   = C.MESON_DIALOG_MESSAGEBOX_TYPE_ERROR
	MessageBoxTypeQuestion                = C.MESON_DIALOG_MESSAGEBOX_TYPE_QUESTION
)

//MessageReceiveHandler is message receive handler.
type MessageReceiveHandler func(int64, string, bool) string

var (
	sendLock       int32
	ReadyChannel   = make(chan struct{})
	requestChannel = make(chan []byte)
	receiveHandler MessageReceiveHandler
)

func LockSendMessage() {
	for {
		if atomic.CompareAndSwapInt32(&sendLock, 0, 1) {
			return
		}
		runtime.Gosched()
	}
}

func TryEnterSendMessage() bool {
	return atomic.CompareAndSwapInt32(&sendLock, 0, 1)
}

func LeaveSendMessage() {
	if !atomic.CompareAndSwapInt32(&sendLock, 1, 0) {
		panic("invalid unlock timing")
	}
}

func MesonFrameworkVersion() string {
	cvers := C.mesonVersions()
	var vers []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&vers))
	bh.Cap = 3
	bh.Len = 3
	bh.Data = uintptr(unsafe.Pointer(cvers))
	ret := fmt.Sprintf("v%d.%d.%d", vers[0], vers[1], vers[2])
	bh.Data = uintptr(0)
	bh.Len = 0
	return ret
}

//export goCallInit
func goCallInit() {
	ReadyChannel <- struct{}{}
}

//export goWaitServerRequest
func goWaitServerRequest() *C.char {
	req := <-requestChannel
	return C.CString(string(req))
}

//export goPostServerResponse
func goPostServerResponse(cid C.uint, crespstr *C.char, needReply C.int) *C.char {
	isReply := needReply > 0
	if isReply {
		LockSendMessage()
		defer LeaveSendMessage()
	}
	id := int64(cid)
	respstr := C.GoString(crespstr)
	resultstr := receiveHandler(id, respstr, isReply)
	var cresultstr *C.char
	if isReply {
		cresultstr = C.CString(resultstr)
	}
	return cresultstr
}

func PostMessage(req []byte) {
	requestChannel <- req
}

func SetMessageReceiveHandler(handler MessageReceiveHandler) {
	receiveHandler = handler
}

func LoadBinding() error {
	fp, err := resolveFrameworkPath()
	if err != nil {
		return err
	}
	cs := C.CString(fp)
	defer C.free(unsafe.Pointer(cs))
	cret := C.loadMesonFramework(cs)
	if cret == C.int(0) {
		return errors.New("load fail framework")
	}

	return nil
}

func RunMesonMainLoop(args []string) int {
	defer C.freeMesonFramework()
	C.MesonApiSetArgc(C.int(len(args)))
	for i, v := range args {
		C.MesonApiAddArgv(C.int(i), C.CString(v))
	}
	C.mesonRegistHandler()
	return int(C.MesonApiMain())
}
