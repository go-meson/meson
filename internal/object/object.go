package object

import (
	"encoding/json"
	"errors"
	"fmt"
	obj "github.com/go-meson/meson/object"
	"log"
	"sync"
)

type CallbackInterface interface {
	Call(obj.ObjectRef, json.RawMessage) (bool, error)
}

type eventRegisters map[int64][]CallbackInterface

type Object struct {
	Id              int64
	ObjType         obj.ObjectType
	registerdEvents eventRegisters
	UserData        interface{}
}

type ObjectRefInternal interface {
	obj.ObjectRef
	EmitEvent(sender ObjectRefInternal, eventID int64, arg json.RawMessage) (bool, error)
}

var (
	lock    = sync.RWMutex{}
	objects = make(map[obj.ObjectType]map[int64]ObjectRefInternal)
)

func NewObject(id int64, objType obj.ObjectType) Object {
	return Object{
		Id:              id,
		ObjType:         objType,
		registerdEvents: make(eventRegisters),
	}
}

func AddObject(t obj.ObjectType, id int64, o ObjectRefInternal) {
	//TODO: need lock?
	lock.Lock()
	defer lock.Unlock()
	tm, ok := objects[t]
	if !ok {
		tm = make(map[int64]ObjectRefInternal)
		objects[t] = tm
	}
	if old, ok := tm[id]; ok {
		panic(fmt.Errorf("object is already exists!(%#v)", old))
	}
	tm[id] = o
}

func GetObject(t obj.ObjectType, id int64) ObjectRefInternal {
	var r ObjectRefInternal
	lock.RLock()
	if tm, ok := objects[t]; ok {
		if o, ok := tm[id]; ok {
			r = o
		}
	}
	lock.RUnlock()
	return r
}

func (o *Object) GetID() int64 {
	return o.Id
}

func (o *Object) GetObjectType() obj.ObjectType {
	return o.ObjType
}

func (o *Object) EmitEvent(sender ObjectRefInternal, eventID int64, args json.RawMessage) (bool, error) {
	var prevent = false

	if events, ok := o.registerdEvents[eventID]; ok {
		for _, e := range events {
			r, err := e.Call(sender, args)
			if err != nil {
				return false, err
			}
			if r {
				prevent = true
			}
		}
	}
	return prevent, nil
}

func (o *Object) Destroy() {
	log.Printf("destroy: %d", o.Id)
}

func (o *Object) Destroyed() {
	log.Printf("destroyed: %d", o.Id)
	lock.Lock()
	defer lock.Unlock()
	if tm, ok := objects[o.ObjType]; ok {
		delete(tm, o.Id)
		if len(tm) == 0 {
			delete(objects, o.ObjType)
		}
	}
}

type objForJSON struct {
	Type obj.ObjectType `json:"type"`
	ID   int64          `json:"id"`
}

func (o *Object) MarshalJSON() ([]byte, error) {
	oj := objForJSON{}
	if o != nil {
		oj.Type = o.ObjType
		oj.ID = o.Id
	}
	return json.Marshal(&oj)
}

func (o *Object) UnmashalJSON(data []byte) error {
	var oj objForJSON
	err := json.Unmarshal(data, &oj)
	if err != nil {
		return err
	}
	obj := GetObject(oj.Type, oj.ID)
	if obj == nil {
		return errors.New("invalid object id")
	}
	*o = *obj.(*Object)
	return nil
}

func (o *Object) AddRegisterdCallback(eventID int64, callback CallbackInterface) int {
	cl := []CallbackInterface{callback}

	var ret int
	if el, ok := o.registerdEvents[eventID]; ok {
		ret = len(el)
		o.registerdEvents[eventID] = append(cl, el...)
	} else {
		o.registerdEvents[eventID] = cl
	}
	return ret
}

func (o *Object) DelRegisterdCallback(eventID int64, no int) bool {
	ret := false
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
			ret = true
		}
	}
	return ret
}

func (o *Object) FindCallback(eventID int64) []CallbackInterface {
	if c, ok := o.registerdEvents[eventID]; ok {
		return c
	}
	return nil
}
