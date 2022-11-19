package event

import (
	"reflect"
)

type PubMode byte

const (
	PubModeSync PubMode = iota
	PubModeAsync
	PubModeAuto
)

type TopicBase interface {
	Name() string
	Hub() Hub
	CurrEventId() EventId
	NewEventId() EventId
	EventType() reflect.Type

	UnSub(name string) bool
	Close(wait bool)
}

type Topic[K any] interface {
	TopicBase

	SubP(name string, lsner Listener[K], qSize uint32) int
	Sub(name string, lsner Listener[K], qSize uint32) (int, error)
	Pub(mode PubMode, evnt K)
}

func NewTopic[K any](name string, hub Hub, evntExample K, logr HubLogger) Topic[K] {
	return NewTopicImpl(name, hub, evntExample, logr)
}
