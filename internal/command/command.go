package command

import (
	"encoding/json"
	"errors"
	"github.com/go-meson/meson/internal/binding"
	"github.com/go-meson/meson/internal/object"
	obj "github.com/go-meson/meson/object"
	"log"
	"sync"
	"sync/atomic"
)

type Command struct {
	Action   binding.ActionType `json:"_action"`
	ActionID int64              `json:"_actionId"`
	Type     obj.ObjectType     `json:"_type"`
	ID       int64              `json:"_id"`
	Method   string             `json:"_method,omitempty"`
	Args     interface{}        `json:"_args"`
}

func MakeCreateCommand(objType obj.ObjectType, args ...interface{}) Command {
	return Command{
		Action: binding.ActCreate,
		Type:   objType,
		Args:   args,
	}
}

func MakeCallCommand(objType obj.ObjectType, id int64, method string, args ...interface{}) Command {
	return Command{
		Action: binding.ActCall,
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

type Response struct {
	Action   binding.ActionType `json:"_action"`
	ActionID int64              `json:"_actionId"`
	Error    string             `json:"_error"`
	ID       int64              `json:"_id"`
	EventID  int64              `json:"_eventId"`
	Method   string             `json:"_method,omitempty"`
	Result   json.RawMessage    `json:"_result"`
}

type ChResp chan *Response

type respHandler func(resp *Response)

var (
	commandID          int64
	APIReady           = false
	responseHandler    = make(map[int64]respHandler)
	requestChannelPool = sync.Pool{New: func() interface{} { return make(ChResp) }}
	commonChannelPool  = sync.Pool{New: func() interface{} { return make(chan interface{}) }}
)

func GetCommonChan() chan interface{} {
	return commonChannelPool.Get().(chan interface{})
}

func ReleaseCommonChan(c chan interface{}) {
	commonChannelPool.Put(c)
}

func getRespChan() ChResp {
	return requestChannelPool.Get().(ChResp)
}

func releaseRespChan(c ChResp) {
	requestChannelPool.Put(c)
}

func MakeRegEventCommand(objType obj.ObjectType, id int64, event string) Command {
	return Command{
		Action: binding.ActRegEvent,
		Type:   objType,
		Method: event,
		ID:     id,
		Args:   &regEventOpt{},
	}
}

func MakeTempEventCommand(objType obj.ObjectType, id int64, numRegist int) Command {
	return Command{
		Action: binding.ActRegEvent,
		Type:   objType,
		ID:     id,
		Args:   &regEventOpt{Temporary: true, Number: numRegist},
	}
}

func MakeUnregEventCommand(objType obj.ObjectType, id int64, eventID int64) Command {
	return Command{
		Action: binding.ActRegEvent,
		Type:   objType,
		ID:     id,
		Args:   &regEventOpt{Delete: true, Number: int(eventID)},
	}
}

func CheckResponse(resp *Response) error {
	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	return nil
}

func PostMessage(cmd *Command) error {
	bytes, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	binding.PostMessage(bytes)
	return nil
}

func sendMessage(cmd *Command, actionID int64, handler respHandler) error {
	cmd.ActionID = actionID
	bytes, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	responseHandler[actionID] = handler
	binding.PostMessage(bytes)
	return nil
}

func SendMessageAsync(cmd *Command, handler respHandler) error {
	actionID := atomic.AddInt64(&commandID, 1)
	return sendMessage(cmd, actionID, handler)
}

func SendMessage(cmd *Command) (json.RawMessage, error) {
	if !binding.TryEnterSendMessage() {
		return nil, errors.New("invalid context")
	}
	defer binding.LeaveSendMessage()
	actionID := atomic.AddInt64(&commandID, 1)
	ch := getRespChan()
	if err := sendMessage(cmd, actionID, func(r *Response) {
		ch <- r
	}); err != nil {
		return nil, err
	}
	resp := <-ch
	releaseRespChan(ch)
	if err := CheckResponse(resp); err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func messageReceived(id int64, msg string, needReply bool) string {
	var resp Response
	err := json.Unmarshal([]byte(msg), &resp)
	if err != nil {
		panic("json decode fail.")
	}
	var result interface{}
	switch resp.Action {
	case binding.ActReply:
		if c, ok := responseHandler[resp.ActionID]; ok {
			delete(responseHandler, resp.ActionID)
			c(&resp)
		} else {
			log.Fatalf("invalid response: %#v\n", resp)
		}
	case binding.ActEvent:
		if o := object.GetObject(id); o == nil {
			log.Fatalf("object not found: %d", id)
		} else {
			if !needReply {
				go func() {
					o.EmitEvent(o, resp.EventID, resp.Result)
				}()
			} else {
				b, _ := o.EmitEvent(o, resp.EventID, resp.Result)
				result = b
			}
		}
	default:
		log.Panicf("invalid action: %#v", resp)
	}

	if needReply {
		r, _ := json.Marshal(result)
		return string(r)
	}
	return ""
}

func init() {
	binding.SetMessageReceiveHandler(messageReceived)
}
