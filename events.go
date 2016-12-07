package meson

import (
	"fmt"
	"github.com/koron/go-dproxy"
	"log"
)

type callbackInterface interface {
	Call(ObjectRef, interface{}) (bool, error)
}

type eventRegister map[int64][]callbackInterface

func (o *object) addCallback(event string, callback callbackInterface) error {
	cmd := makeRegEventCommand(o.objType, o.id, event)

	resp, err := sendMessage(&cmd)
	if err != nil {
		return err
	}

	eventID, err := dproxy.New(resp).Int64()
	if err != nil {
		return err
	}

	o.addRegisterdCallback(eventID, callback)

	return nil
}

type tempEventItem struct {
	id   int64
	name string
}

func (o *object) makeTempEventAsync(num int, handler func([]tempEventItem, error)) {
	cmd := makeTempEventCommand(o.objType, o.id, num)
	sendMessageAsync(&cmd, func(r *response) {
		if err := checkResponse(r); err != nil {
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
		ret := make([]tempEventItem, resplen)
		for idx := 0; idx < resplen; idx++ {
			ret[idx].id, _ = rproxy.A(idx).M("eventId").Int64()
			ret[idx].name, _ = rproxy.A(idx).M("eventName").String()
		}
		handler(ret, nil)
	})
}

func (o *object) makeTemporaryEvents(num int) ([]tempEventItem, error) {
	ch := getCommonChan()
	o.makeTempEventAsync(num, func(items []tempEventItem, err error) {
		if err != nil {
			ch <- err
		} else {
			ch <- items
		}
	})
	ret := <-ch
	releaseCommonChan(ch)
	if err, ok := ret.(error); ok {
		return nil, err
	}
	return ret.([]tempEventItem), nil
}

func (o *object) addRegisterdCallback(eventID int64, callback callbackInterface) int {
	cl := []callbackInterface{callback}

	var ret int
	if el, ok := o.registerdEvents[eventID]; ok {
		ret = len(el)
		o.registerdEvents[eventID] = append(cl, el...)
	} else {
		o.registerdEvents[eventID] = cl
	}
	return ret
}

func (o *object) delRegisterdCallback(eventID int64, no int) {
	if el, ok := o.registerdEvents[eventID]; ok {
		if no == 0 {
			el = nil
		} else {
			el = append(el[:no], el[no+1:]...)
		}
		if len(el) > 0 {
			o.registerdEvents[eventID] = el
		} else {
			delete(o.registerdEvents, eventID)
			cmd := makeUnregEventCommand(o.objType, o.id, eventID)
			postMessage(&cmd)
		}
	}
}

func (o *object) findCallback(eventID int64) []callbackInterface {
	if c, ok := o.registerdEvents[eventID]; ok {
		return c
	}
	return nil
}

type CommonCallbackHandler func(ObjectRef)

type commonCallbackItem struct {
	f CommonCallbackHandler
}

func (p commonCallbackItem) Call(o ObjectRef, arg interface{}) (bool, error) {
	args, ok := arg.([]interface{})
	if !ok || len(args) != 0 {
		log.Panicf("Invalid arg type: %#v", arg)
	}
	p.f(o)
	return false, nil
}

type CommonPreventableCallbackHandler func(ObjectRef) bool

type commonPreventableCallbackItem struct {
	f CommonPreventableCallbackHandler
}

func (p commonPreventableCallbackItem) Call(o ObjectRef, arg interface{}) (bool, error) {
	args, ok := arg.([]interface{})
	if !ok || len(args) != 0 {
		log.Panicf("Invalid arg type: %#v", arg)
	}
	return p.f(o), nil
}
