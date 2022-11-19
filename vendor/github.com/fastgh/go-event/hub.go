package event

import (
	"fmt"
	"reflect"
	"sync"
)

type HubT struct {
	name   string
	mx     sync.RWMutex
	topics map[string]TopicBase
	logr   HubLogger
}

type Hub = *HubT

func (me Hub) Logger() HubLogger { return me.logr }

func (me Hub) Name() string { return me.name }

func (me Hub) RegisterTopic(topic TopicBase) {
	me.mx.Lock()
	defer me.mx.Unlock()

	nm := topic.Name()
	me.logr.LogInfo(TopicRegisterBegin, nm, "")

	if _, has := me.topics[nm]; has {
		panic(fmt.Errorf("duplicated topic '%s'", nm))
	}

	me.topics[nm] = topic

	me.logr.LogInfo(TopicRegisterOk, nm, "")
}

func (me Hub) HasTopic(name string) bool {
	me.mx.RLock()
	defer me.mx.RUnlock()

	_, has := me.topics[name]
	return has
}

func (me Hub) GetTopic(name string, evntExample any) TopicBase {
	me.mx.RLock()
	defer me.mx.RUnlock()

	r, has := me.topics[name]
	if has {
		expTyp := reflect.TypeOf(evntExample)
		actualTyp := r.EventType()
		if expTyp != actualTyp {
			panic(fmt.Errorf("expected event type is %v, but got %v", expTyp, actualTyp))
		}
	}

	return r
}

func (me Hub) Close(wait bool) {
	me.mx.RLock()
	defer me.mx.RUnlock()

	me.logr.LogInfo(HubCloseBegin, "", "")

	for _, tp := range me.topics {
		tp.Close(wait)
	}

	me.logr.LogInfo(HubCloseOk, "", "")
}
