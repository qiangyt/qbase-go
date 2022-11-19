package event

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

type TopicImpl[K any] struct {
	name   string
	hub    Hub
	typ    reflect.Type
	lsners []*EventListener[K]
	waitG  sync.WaitGroup
	logr   TopicLogger
	eid    atomic.Int64
	mx     sync.RWMutex
}

func NewTopicImpl[K any](name string, hub Hub, example K, logr HubLogger) *TopicImpl[K] {
	return &TopicImpl[K]{
		name:   name,
		hub:    hub,
		typ:    reflect.TypeOf(example),
		lsners: []*EventListener[K]{},
		waitG:  sync.WaitGroup{},
		eid:    atomic.Int64{},
		mx:     sync.RWMutex{},
		logr:   NewTopicLogger(name, logr),
	}
}

func (me *TopicImpl[K]) EventType() reflect.Type { return me.typ }

func (me *TopicImpl[K]) Name() string { return me.name }

func (me *TopicImpl[K]) Hub() Hub { return me.hub }

func (me *TopicImpl[K]) CurrEventId() EventId { return EventId(me.eid.Load()) }

func (me *TopicImpl[K]) NewEventId() EventId { return EventId(me.eid.Add(1)) }

func (me *TopicImpl[K]) SubP(name string, lsner Listener[K], qSize uint32) int {
	r, err := me.Sub(name, lsner, qSize)
	if err != nil {
		panic(err)
	}
	return r
}

func (me *TopicImpl[K]) Sub(name string, lsner Listener[K], qSize uint32) (int, error) {
	if lsner == nil {
		return 0, errors.New("listener cannot be nil")
	}

	me.mx.Lock()
	defer me.mx.Unlock()

	evntLsners := me.lsners
	for i, existing := range evntLsners {
		if existing.name == name {
			me.logr.LogError(ListenerSubErr, name, fmt.Sprintf("duplicated listener on #%d", i))
			return -1, nil
		}
	}

	evntLsner := NewEventListener(name, lsner, qSize, me.logr)
	evntLsner.Start()

	evntLsners = append(evntLsners, evntLsner)
	me.lsners = evntLsners
	me.waitG.Add(1)

	me.logr.LogInfo(ListenerSubOk, name)
	return len(evntLsners), nil
}

func (me *TopicImpl[K]) UnSub(name string) bool {
	me.mx.Lock()
	defer me.mx.Unlock()

	lsners := me.lsners
	for i, existing := range lsners {
		if existing.name == name {
			me.lsners = append(lsners[:i], lsners[i+1])
			me.logr.LogInfo(ListenerUnsubOk, name)

			stopEvent := NewCloseEvent(me.NewEventId(), me.Hub().Name(), me.name)
			me.stopListener(existing, stopEvent)
			return true
		}
	}

	me.logr.LogError(ListenerUnsubErr, name, "not found")
	return false
}

func (me *TopicImpl[K]) Pub(mode PubMode, evnt K) {
	var async bool
	if mode == PubModeAsync {
		async = true
	} else if mode == PubModeSync {
		async = false
	} else {
		async = len(me.lsners) >= 100
	}

	if async {
		go me.doPub(evnt)
	} else {
		me.doPub(evnt)
	}
}

func (me *TopicImpl[K]) doPub(evntData K) {
	evnt := NewDataEvent(me.NewEventId(), me.Hub().name, me.name, evntData)

	me.mx.RLock()
	defer me.mx.RUnlock()

	me.logr.LogEventDebug(EventPubBegin, "", evnt)

	for _, lsner := range me.lsners {
		lsner.SendEvent(evnt)
	}

	me.logr.LogEventDebug(EventPubOk, "", evnt)
}

func (me *TopicImpl[K]) Close(wait bool) {
	me.mx.RLock()
	defer me.mx.RUnlock()

	stopEvnt := NewCloseEvent(me.NewEventId(), me.Hub().Name(), me.name)
	me.logr.LogEventDebug(TopicCloseBegin, "", stopEvnt)

	for _, lsner := range me.lsners {
		me.stopListener(lsner, stopEvnt)
	}

	if wait {
		me.waitG.Wait()
	}

	me.logr.LogEventDebug(TopicCloseOk, "", stopEvnt)
}

func (me *TopicImpl[K]) stopListener(lsner *EventListener[K], stopEvnt Event) {
	lsner.Stop(stopEvnt)
	me.waitG.Done()
}
