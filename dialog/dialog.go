package dialog

import (
	"errors"
	"github.com/go-meson/meson/app"
	"github.com/go-meson/meson/internal/binding"
	"github.com/go-meson/meson/internal/command"
	"github.com/go-meson/meson/internal/event"
	obj "github.com/go-meson/meson/internal/object"
	"github.com/go-meson/meson/object"
	"github.com/go-meson/meson/window"
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

type MessageBoxType binding.MessageBoxType

type msgBoxOpt struct {
	Type    MessageBoxType `json:"type"`
	Title   string         `json:"title"`
	Message string         `json:"message"`
	MessageBoxOpt
}

const (
	MessageBoxTypeNone     MessageBoxType = MessageBoxType(binding.MessageBoxTypeNone)
	MessageBoxTypeInfo                    = MessageBoxType(binding.MessageBoxTypeInfo)
	MessageBoxTypeWarning                 = MessageBoxType(binding.MessageBoxTypeWarning)
	MessageBoxTypeError                   = MessageBoxType(binding.MessageBoxTypeError)
	MessageBoxTypeQuestion                = MessageBoxType(binding.MessageBoxTypeQuestion)
)

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

func ShowMessageBox(window *window.Window, message string, title string, messageBoxType MessageBoxType, opt *MessageBoxOpt) (int, error) {
	tmpl := makeMsgBoxOpt(message, title, messageBoxType, opt)
	var winid int64
	if window != nil {
		winid = window.Id
	}
	cmd := command.MakeCallCommand(binding.ObjApp, binding.ObjAppID, "showMessageBox", winid, &tmpl)
	r, err := command.SendMessage(&cmd)
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

func (mb msgBoxCallbackItem) Call(o object.ObjectRef, arg interface{}) (bool, error) {
	app := o.(*app.App)
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
	event.DeleteRegisterdCallback(&app.Object, mb.eventID, mb.eventNo)
	return false, nil
}

func ShowMessageBoxAsync(window *window.Window, message string, title string, messageBoxType MessageBoxType, opt *MessageBoxOpt, handler func(int, error)) {
	if handler == nil {
		panic(errors.New("invalid argument"))
	}
	tmpl := makeMsgBoxOpt(message, title, messageBoxType, opt)
	var winid int64
	if window != nil {
		winid = window.Id
	}

	app := obj.GetObject(binding.ObjAppID).(*app.App)
	event.MakeTempEventAsync(&app.Object, 1, func(items []event.TempEventItem, err error) {
		if err != nil {
			handler(-1, err)
			return
		}
		eventID := items[0].Id
		eventName := items[0].Name
		item := &msgBoxCallbackItem{f: handler, eventID: eventID}
		eventNo := app.AddRegisterdCallback(eventID, item)
		item.eventNo = eventNo
		cmd := command.MakeCallCommand(app.ObjType, app.Id, "showMessageBox", winid, &tmpl, eventName)
		command.SendMessageAsync(&cmd, func(r *command.Response) {
			if err := command.CheckResponse(r); err != nil {
				handler(-1, err)
				event.DeleteRegisterdCallback(&app.Object, eventID, eventNo)
				return
			}
		})
	})

	cmd := command.MakeTempEventCommand(app.ObjType, app.Id, 1)
	command.SendMessageAsync(&cmd, func(r *command.Response) {
		if err := command.CheckResponse(r); err != nil {
			handler(-1, err)
			return
		}

	})
}
