package meson

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type object struct {
	id              int64
	objType         objectType
	registerdEvents eventRegister
	UserData        interface{}
}

type ObjectRef interface {
	getID() int64
	getObjectType() objectType
	emitEvent(sender ObjectRef, eventID int64, arg interface{}) (bool, error)
}

func newObject(id int64, objType objectType) object {
	o := object{
		id:              id,
		objType:         objType,
		registerdEvents: make(eventRegister),
	}
	return o
}

var objects = make(map[int64]ObjectRef)

func addObject(id int64, o ObjectRef) {
	if old, ok := objects[id]; ok {
		panic(fmt.Errorf("object is already exists!(%#v)", old))
	}
	objects[id] = o
}

func getObject(id int64) ObjectRef {
	if o, ok := objects[id]; ok {
		return o
	}
	return nil
}

func (o *object) getID() int64 {
	return o.id
}

func (o *object) getObjectType() objectType {
	return o.objType
}

func (o *object) emitEvent(sender ObjectRef, eventID int64, arg interface{}) (bool, error) {
	var prevent = false

	if events, ok := o.registerdEvents[eventID]; ok {
		for _, e := range events {
			r, err := e.Call(sender, arg)
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

func (o *object) Destroy() {
	log.Printf("destroy: %d", o.id)
}

func (o *object) Destroyed() {
	log.Printf("destroyed: %d", o.id)
	delete(objects, o.id)
}

func (o *object) MarshalJSON() ([]byte, error) {
	var id int64
	if o != nil {
		id = o.id
	}
	return json.Marshal(&id)
}

func (o *object) UnmashalJSON(data []byte) error {
	var id int64
	err := json.Unmarshal(data, &id)
	if err != nil {
		return err
	}
	obj := getObject(id)
	if obj == nil {
		return errors.New("Invalid Object ID")
	}
	*o = *obj.(*object)
	return nil
}
