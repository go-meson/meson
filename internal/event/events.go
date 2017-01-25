package event

import (
	"fmt"
	evt "github.com/go-meson/meson/event"
	"github.com/go-meson/meson/internal/command"
	"github.com/go-meson/meson/internal/object"
	obj "github.com/go-meson/meson/object"
	"github.com/koron/go-dproxy"
	"log"
)

func AddCallback(o *object.Object, event string, callback object.CallbackInterface) error {
	cmd := command.MakeRegEventCommand(o.ObjType, o.Id, event)

	resp, err := command.SendMessage(&cmd)
	if err != nil {
		return err
	}

	eventID, err := dproxy.New(resp).Int64()
	if err != nil {
		return err
	}

	o.AddRegisterdCallback(eventID, callback)

	return nil
}

type TempEventItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func MakeTempEventAsync(o *object.Object, num int, handler func([]TempEventItem, error)) {
	cmd := command.MakeTempEventCommand(o.ObjType, o.Id, num)
	command.SendMessageAsync(&cmd, func(r *command.Response) {
		if err := command.CheckResponse(r); err != nil {
			handler(nil, err)
			return
		}
		rproxy := dproxy.New(r.Result)
		resplen := 0
		if rarray, err := rproxy.Array(); err == nil {
			resplen = len(rarray)
		} else {
			handler(nil, err)
			return
		}
		if num != resplen {
			handler(nil, fmt.Errorf("Response length error %d -> %d", num, resplen))
			return
		}
		ret := make([]TempEventItem, resplen)
		for idx := 0; idx < resplen; idx++ {
			ret[idx].Id, _ = rproxy.A(idx).M("eventId").Int64()
			ret[idx].Name, _ = rproxy.A(idx).M("eventName").String()
		}
		handler(ret, nil)
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

func (p CommonCallbackItem) Call(o obj.ObjectRef, arg interface{}) (bool, error) {
	args, ok := arg.([]interface{})
	if !ok || len(args) != 0 {
		log.Panicf("Invalid arg type: %#v", arg)
	}
	p.F(o)
	return false, nil
}

type CommonPreventableCallbackItem struct {
	F evt.CommonPreventableCallbackHandler
}

func (p CommonPreventableCallbackItem) Call(o obj.ObjectRef, arg interface{}) (bool, error) {
	args, ok := arg.([]interface{})
	if !ok || len(args) != 0 {
		log.Panicf("Invalid arg type: %#v", arg)
	}
	return p.F(o), nil
}
