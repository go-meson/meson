package meson

import (
	"errors"
)

/*
type actionType string

const (
	actCreate = "create"
	actDelete = "delete"
	actCall   = "call"
	actReply  = "reply"
	actEvent  = "event"
)
*/

type command struct {
	Action   actionType  `json:"_action"`
	ActionID int64       `json:"_actionId"`
	Type     objectType  `json:"_type"`
	ID       int64       `json:"_id"`
	Method   string      `json:"_method,omitempty"`
	Args     interface{} `json:"_args"`
}

func makeCreateCommand(objType objectType, args ...interface{}) command {
	return command{
		Action: actCreate,
		Type:   objType,
		Args:   args,
	}
}

func makeCallCommand(objType objectType, id int64, method string, args ...interface{}) command {
	return command{
		Action: actCall,
		Type:   objType,
		Method: method,
		ID:     id,
		Args:   args,
	}
}

type regEventOpt struct {
	Delete    bool `json:"delete"`
	Temporary bool `json:"temporary"`
	Number    int  `json:"number"`
}

func makeRegEventCommand(objType objectType, id int64, event string) command {
	return command{
		Action: actRegEvent,
		Type:   objType,
		Method: event,
		ID:     id,
		Args:   &regEventOpt{},
	}
}

func makeTempEventCommand(objType objectType, id int64, numRegist int) command {
	return command{
		Action: actRegEvent,
		Type:   objType,
		ID:     id,
		Args:   &regEventOpt{Temporary: true, Number: numRegist},
	}
}

func makeUnregEventCommand(objType objectType, id int64, eventID int64) command {
	return command{
		Action: actRegEvent,
		Type:   objType,
		ID:     id,
		Args:   &regEventOpt{Delete: true, Number: int(eventID)},
	}
}

type response struct {
	Action   actionType  `json:"_action"`
	ActionID int64       `json:"_actionId"`
	Error    string      `json:"_error"`
	ID       int64       `json:"_id"`
	EventID  int64       `json:"_eventId"`
	Method   string      `json:"_method,omitempty"`
	Result   interface{} `json:"_result"`
}

func checkResponse(resp *response) error {
	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	return nil
}
