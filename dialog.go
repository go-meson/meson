package meson

import (
	"errors"
	"github.com/koron/go-dproxy"
	"log"
)

type MessageBoxOpt struct {
	Buttons   []string `json:"buttons"`
	DefaultID int      `json:"defaultId"`
	CancelID  int      `json:"cancelId"`
	Detail    string   `json:"detail"`
	NoLink    bool     `json:"noLink"`
}

type msgBoxOpt struct {
	Type    MessageBoxType `json:"type"`
	Title   string         `json:"title"`
	Message string         `json:"message"`
	MessageBoxOpt
}

func makeMsgBoxOpt(message string, title string, messageBoxType MessageBoxType, opt *MessageBoxOpt) msgBoxOpt {
	tmpl := msgBoxOpt{Type: messageBoxType, Title: title, Message: message}
	if opt != nil {
		tmpl.MessageBoxOpt = *opt
	}
	if tmpl.Buttons == nil {
		switch tmpl.Type {
		case MessageBoxTypeNone, MessageBoxTypeInfo, MessageBoxTypeError, MessageBoxTypeWarning:
			tmpl.Buttons = []string{"OK"}
		case MessageBoxTypeQuestion:
			tmpl.Buttons = []string{"YES", "NO"}
		}
	}
	return tmpl
}

func ShowMessageBox(window *Window, message string, title string, messageBoxType MessageBoxType, opt *MessageBoxOpt) (int, error) {
	tmpl := makeMsgBoxOpt(message, title, messageBoxType, opt)
	var winid int64
	if window != nil {
		winid = window.id
	}
	cmd := makeCallCommand(objApp, objAppID, "showMessageBox", winid, &tmpl)
	r, err := sendMessage(&cmd)
	if err != nil {
		return -1, err
	}
	buttonID, err := dproxy.New(r).Int64()
	if err != nil {
		return -1, err
	}
	return int(buttonID), nil
}

type msgBoxCallbackHandler func(int, error)

type msgBoxCallbackItem struct {
	f       msgBoxCallbackHandler
	eventID int64
	eventNo int
}

func (mb msgBoxCallbackItem) Call(o ObjectRef, arg interface{}) (bool, error) {
	app := o.(*App)
	args, ok := arg.([]interface{})
	if !ok || len(args) != 1 {
		log.Panicf("Invalid arg type: %#v", arg)
	}
	button, err := dproxy.New(args[0]).Int64()
	if err != nil {
		mb.f(-1, err)
	} else {
		mb.f(int(button), nil)
	}
	log.Printf("eventID = %d, eventNo = %d\n", mb.eventID, mb.eventNo)
	app.delRegisterdCallback(mb.eventID, mb.eventNo)
	return false, nil
}

func ShowMessageBoxAsync(window *Window, message string, title string, messageBoxType MessageBoxType, opt *MessageBoxOpt, handler func(int, error)) {
	if handler == nil {
		panic(errors.New("invalid argument"))
	}
	tmpl := makeMsgBoxOpt(message, title, messageBoxType, opt)
	var winid int64
	if window != nil {
		winid = window.id
	}

	app := getObject(objAppID).(*App)
	app.makeTempEventAsync(1, func(items []tempEventItem, err error) {
		if err != nil {
			handler(-1, err)
			return
		}
		eventID := items[0].id
		eventName := items[0].name
		item := &msgBoxCallbackItem{f: handler, eventID: eventID}
		eventNo := app.addRegisterdCallback(eventID, item)
		item.eventNo = eventNo
		cmd := makeCallCommand(app.objType, app.id, "showMessageBox", winid, &tmpl, eventName)
		sendMessageAsync(&cmd, func(r *response) {
			if err := checkResponse(r); err != nil {
				handler(-1, err)
				app.delRegisterdCallback(eventID, eventNo)
				return
			}
		})
	})
	cmd := makeTempEventCommand(objApp, objAppID, 1)
	sendMessageAsync(&cmd, func(r *response) {
		if err := checkResponse(r); err != nil {
			handler(-1, err)
			return
		}

	})
}
