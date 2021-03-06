package dialog

import (
	"encoding/json"
	"errors"
	"github.com/go-meson/meson/internal/binding"
	"github.com/go-meson/meson/internal/command"
	"github.com/go-meson/meson/internal/event"
	obj "github.com/go-meson/meson/internal/object"
	"github.com/go-meson/meson/object"
	"github.com/go-meson/meson/window"
	"log"
)

// MessageBoxOpt is dialog box creation options.
type MessageBoxOpt struct {
	Buttons   []string `json:"buttons"`   // Array of texts for buttons.
	DefaultID int      `json:"defaultId"` // Index in the Buttons array which will be selected by default when the message box opened.
	CancelID  int      `json:"cancelId"`  // Index in the Buttons array which will be selected when user cancels the dialog instead of clicking the buttons of the dialog.
	Detail    string   `json:"detail"`    // Extra information of the message
	NoLink    bool     `json:"noLink"`    // TODO:
}

type msgBoxOpt struct {
	Type    MessageBoxType `json:"type"`
	Title   string         `json:"title"`
	Message string         `json:"message"`
	MessageBoxOpt
}

// MessageBoxType is following message box type.
type MessageBoxType binding.MessageBoxType

const (
	MessageBoxTypeNone     MessageBoxType = MessageBoxType(binding.MessageBoxTypeNone)     // without icons
	MessageBoxTypeInfo                    = MessageBoxType(binding.MessageBoxTypeInfo)     // with icon for message information.
	MessageBoxTypeWarning                 = MessageBoxType(binding.MessageBoxTypeWarning)  // with icon for message warning.
	MessageBoxTypeError                   = MessageBoxType(binding.MessageBoxTypeError)    // with icon for message error.
	MessageBoxTypeQuestion                = MessageBoxType(binding.MessageBoxTypeQuestion) // with icon for question message.
)

var (
	dlgCls = func() *obj.Object {
		dlg := obj.NewObject(binding.ObjStaticID, binding.ObjDialog)
		obj.AddObject(binding.ObjDialog, binding.ObjStaticID, &dlg)
		return &dlg
	}()
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
	//TODO: dialog static method
	cmd := command.MakeCallCommand(binding.ObjDialog, binding.ObjStaticID, "showMessageBox", winid, &tmpl)
	r, err := command.SendMessage(&cmd)
	if err != nil {
		return -1, err
	}
	var buttonID int
	err = json.Unmarshal(r, &buttonID)
	if err != nil {
		return -1, err
	}
	return buttonID, nil
}

type msgBoxCallbackHandler func(int, error)

type msgBoxCallbackItem struct {
	f       msgBoxCallbackHandler
	eventID int64
	eventNo int
}

func (mb msgBoxCallbackItem) Call(o object.ObjectRef, arg json.RawMessage) (bool, error) {
	args := struct {
		ButtonID int `json:"buttonID"`
	}{}
	err := json.Unmarshal(arg, &args)
	if err != nil {
		mb.f(-1, err)
	} else {
		mb.f(args.ButtonID, nil)
	}
	obj := o.(*obj.Object)
	log.Printf("eventID = %d, eventNo = %d\n", mb.eventID, mb.eventNo)
	event.DeleteRegisterdCallback(obj, mb.eventID, mb.eventNo)
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

	event.MakeTempEventAsync(dlgCls, 1, func(items []event.TempEventItem, err error) {
		if err != nil {
			handler(-1, err)
			return
		}
		eventID := items[0].EventID
		eventName := items[0].EventName
		item := &msgBoxCallbackItem{f: handler, eventID: eventID}
		eventNo := dlgCls.AddRegisterdCallback(eventID, item)
		item.eventNo = eventNo
		cmd := command.MakeCallCommand(dlgCls.ObjType, dlgCls.Id, "showMessageBox", winid, &tmpl, eventName)
		if err := command.SendMessageAsync(&cmd, func(r *command.Response) {
			if err := command.CheckResponse(r); err != nil {
				handler(-1, err)
				event.DeleteRegisterdCallback(dlgCls, eventID, eventNo)
				return
			}
		}); err != nil {
			panic(err)
		}
	})

	cmd := command.MakeTempEventCommand(dlgCls.ObjType, dlgCls.Id, 1)
	if err := command.SendMessageAsync(&cmd, func(r *command.Response) {
		if err := command.CheckResponse(r); err != nil {
			handler(-1, err)
			return
		}

	}); err != nil {
		panic(err)
	}
}
