package meson

/*
#cgo CFLAGS:	-mmacosx-version-min=10.9
#cgo CFLAGS:	-DMESON_DIR=${SRCDIR}
#cgo CFLAGS: -Idist/include
#cgo LDFLAGS: -F${SRCDIR}/dist
#cgo LDFLAGS: -framework Meson
#cgo LDFLAGS: -Wl,-rpath,@executable_path/../Frameworks
#cgo LDFLAGS: -Wl,-rpath,${SRCDIR}/dist
#cgo LDFLAGS: -mmacosx-version-min=10.9

#include <stdlib.h>
#include "meson.h"
extern void mesonRegistHandler(void);
*/
import "C"
import "encoding/json"
import "log"
import "sync"

type actionType int

const (
	actCreate   actionType = C.MESON_ACTION_TYPE_CREATE
	actDelete              = C.MESON_ACTION_TYPE_DELETE
	actCall                = C.MESON_ACTION_TYPE_CALL
	actReply               = C.MESON_ACTION_TYPE_REPLY
	actEvent               = C.MESON_ACTION_TYPE_EVENT
	actRegEvent            = C.MESON_ACTION_TYPE_REGISTER_EVENT
)

const (
	objAppID int64 = C.MESON_OBJID_APP
)

type objectType int

const (
	objApp         objectType = C.MESON_OBJECT_TYPE_APP
	objWindow                 = C.MESON_OBJECT_TYPE_WINDOW
	objSession                = C.MESON_OBJECT_TYPE_SESSION
	objWebContents            = C.MESON_OBJECT_TYPE_WEB_CONTENTS
	objMenu                   = C.MESON_OBJECT_TYPE_MENU
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

type chResp chan *response

type respHandler func(resp *response)

var (
	readyChannel       = make(chan struct{})
	requestChannel     = make(chan *command)
	responseHandler    = make(map[int64]respHandler)
	requestChannelPool = sync.Pool{New: func() interface{} { return make(chResp) }}
	commonChannelPool  = sync.Pool{New: func() interface{} { return make(chan interface{}) }}
)

func getCommonChan() chan interface{} {
	return commonChannelPool.Get().(chan interface{})
}

func releaseCommonChan(c chan interface{}) {
	commonChannelPool.Put(c)
}

func getRespChan() chResp {
	return requestChannelPool.Get().(chResp)
}
func releaseRespChan(c chResp) {
	requestChannelPool.Put(c)
}

func postMessage(cmd *command) {
	requestChannel <- cmd
}

func sendMessageInternal(cmd *command, actionID int64, handler respHandler) {
	responseHandler[actionID] = handler
	cmd.ActionID = actionID
	requestChannel <- cmd
}

//export goCallInit
func goCallInit() {
	readyChannel <- struct{}{}
}

//export goWaitServerRequest
func goWaitServerRequest() *C.char {
	req := <-requestChannel
	str, err := json.Marshal(req)
	if err != nil {
		panic("cmd encode error!e")
	}
	log.Printf("S: %s\n", str)
	return C.CString(string(str))
}

//export goPostServerResponse
func goPostServerResponse(cid C.uint, crespstr *C.char, needReply C.int) *C.char {
	var resultstr *C.char
	var result interface{}
	if needReply > 0 {
		lockSendMessage()
		defer leaveSendMessage()
	}

	id := int64(cid)
	respstr := C.GoString(crespstr)
	var resp response
	err := json.Unmarshal([]byte(respstr), &resp)
	if err != nil {
		panic("json decode fail.")
	}
	switch resp.Action {
	case actReply:
		log.Printf("R: %#v\n", resp)
		if c, ok := responseHandler[resp.ActionID]; ok {
			delete(responseHandler, resp.ActionID)
			c(&resp)
		} else {
			log.Fatalf("invalid response: %#v\n", resp)
		}
	case actEvent:
		log.Printf("E: %#v\n", resp)
		if o := getObject(id); o == nil {
			log.Fatalf("object not found %d\n", id)
		} else {
			if needReply == 0 {
				go func() {
					o.emitEvent(o, resp.EventID, resp.Result)
				}()
			} else {
				b, _ := o.emitEvent(o, resp.EventID, resp.Result)
				result = b
			}
		}
	default:
		log.Panicf("invalid action: %#v", resp)
	}
	if needReply != 0 {
		r, _ := json.Marshal(result)
		resultstr = C.CString(string(r))
	}
	return resultstr
}

func runMesonMainLoop(args []string) int {
	C.MesonApiSetArgc(C.int(len(args)))
	for i, v := range args {
		C.MesonApiAddArgv(C.int(i), C.CString(v))
	}
	C.mesonRegistHandler()
	return int(C.MesonApiMain())
}
