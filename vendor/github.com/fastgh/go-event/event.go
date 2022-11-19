package event

import (
	"encoding/json"
	"fmt"

	phuslu "github.com/phuslu/log"
	"github.com/pkg/errors"
)

type EventId int64

type EventT struct {
	Id    EventId
	Hub   string
	Topic string
	Close bool
	Data  any
}

type Event = *EventT

func NewDataEvent(id EventId, hub string, topic string, dat any) Event {
	return &EventT{
		Id:    id,
		Hub:   hub,
		Topic: topic,
		Data:  dat,
		Close: false,
	}
}

func NewCloseEvent(id EventId, hub string, topic string) Event {
	return &EventT{
		Id:    id,
		Hub:   hub,
		Topic: topic,
		Data:  nil,
		Close: true,
	}
}

func (me Event) String() string {
	bytes, err := json.Marshal(me)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal event"))
	}
	return string(bytes)
}

func (me Event) MarshalObject(entry *phuslu.Entry) {
	entry.Int64("id", int64(me.Id)).
		Str("hub", me.Hub).
		Str("topic", me.Topic).
		Bool("close", me.Close)

	dat := me.Data
	if dat == nil {
		return
	}

	datStr, is := dat.(string)
	if is {
		entry.Str("data", datStr)
		return
	}

	entry.Str("data", fmt.Sprintf("%v", dat))
}
