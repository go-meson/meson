package event

import (
	"encoding/json"
	evt "github.com/go-meson/meson/event"
	"github.com/go-meson/meson/internal/command"
	"github.com/go-meson/meson/internal/object"
	obj "github.com/go-meson/meson/object"
)

func AddCallback(o *object.Object, event string, callback object.CallbackInterface) error {
	cmd := command.MakeRegEventCommand(o.ObjType, o.Id, event)

	resp, err := command.SendMessage(&cmd)
	if err != nil {
		return err
	}

	var eventID int64
	err = json.Unmarshal(resp, &eventID)
	if err != nil {
		return err
	}

	o.AddRegisterdCallback(eventID, callback)

	return nil
}

type TempEventItem struct {
	EventID   int64  `json:"eventId"`
	EventName string `json:"eventName"`
}

func MakeTempEventAsync(o *object.Object, num int, handler func([]TempEventItem, error)) {
	cmd := command.MakeTempEventCommand(o.ObjType, o.Id, num)
	command.SendMessageAsync(&cmd, func(r *command.Response) {
		if err := command.CheckResponse(r); err != nil {
			handler(nil, err)
			return
		}
		var events []TempEventItem
		if err := json.Unmarshal(r.Result, &events); err != nil {
			handler(nil, err)
			return
		}
		handler(events, nil)
	})
}

func MakeTemporaryEvents(o *object.Object, num int) ([]TempEventItem, error) {
	ch := command.GetCommonChan()
	MakeTempEventAsync(o, num, func(items []TempEventItem, err error) {
		if err != nil {
			ch <- err
		} else {
			ch <- items
		}
	})
	ret := <-ch
	command.ReleaseCommonChan(ch)
	if err, ok := ret.(error); ok {
		return nil, err
	}
	return ret.([]TempEventItem), nil
}

func DeleteRegisterdCallback(o *object.Object, eventID int64, no int) {
	if o.DelRegisterdCallback(eventID, no) {
		cmd := command.MakeUnregEventCommand(o.ObjType, o.Id, eventID)
		command.PostMessage(&cmd)
	}

}

type CommonCallbackItem struct {
	F evt.CommonCallbackHandler
}

func (p CommonCallbackItem) Call(o obj.ObjectRef, arg json.RawMessage) (bool, error) {
	p.F(o)
	return false, nil
}

type CommonPreventableCallbackItem struct {
	F evt.CommonPreventableCallbackHandler
}

func (p CommonPreventableCallbackItem) Call(o obj.ObjectRef, arg json.RawMessage) (bool, error) {
	return p.F(o), nil
}
