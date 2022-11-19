package event

import "sync"

func NewHub(name string, logr Logger) Hub {
	return &HubT{
		name:   name,
		mx:     sync.RWMutex{},
		topics: map[string]TopicBase{},
		logr:   NewHubLogger(name, logr),
	}
}

func CreateTopic[K any](hub Hub, name string, evntExample K) Topic[K] {
	r := NewTopic(name, hub, evntExample, hub.Logger())
	hub.RegisterTopic(r)
	return r
}

func GetTopic[K any](me Hub, name string, evntExample K) Topic[K] {
	r := me.GetTopic(name, evntExample)
	return r.(Topic[K])
}
