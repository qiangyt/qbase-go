package qplugin

import (
	"fmt"
	"sync"

	"github.com/fastgh/go-comm/v2"
)

type BasePluginLoaderT struct {
	started   bool
	namespace string
	plugins   map[string]Plugin
	mutex     sync.RWMutex
}

type BasePluginLoader = *BasePluginLoaderT

func NewPluginLoader(namespace string) BasePluginLoader {
	return &BasePluginLoaderT{
		started:   false,
		namespace: namespace,
		plugins:   map[string]Plugin{},
		mutex:     sync.RWMutex{},
	}
}

func (me BasePluginLoader) Register(plugin Plugin) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	name := plugin.Name()
	if _, found := me.plugins[name]; found {
		panic(fmt.Errorf("plugin %s is duplicated", PluginId(me.Namespace(), name)))
	}
	me.plugins[name] = plugin
}

func (me BasePluginLoader) RegisterThenStart(logger comm.Logger, plugin Plugin) {
	me.Register(plugin)

	if err := StartPlugin(me.Namespace(), plugin, logger); err != nil {
		panic(err)
	}
}

func (me BasePluginLoader) Namespace() string {
	return me.namespace
}

func (me BasePluginLoader) Plugins() map[string]Plugin {
	return me.plugins
}

func (me BasePluginLoader) Start(logger comm.Logger) error {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	if me.started {
		logger.Info().Msg("started, already")
		return nil
	}

	errs := comm.NewErrorGroup(false)
	ns := me.Namespace()

	for _, plugin := range me.plugins {
		if err := StartPlugin(ns, plugin, logger); err != nil {
			errs.Add(err)
		}
	}

	if errs.HasError() {
		return errs
	}

	me.started = true
	return nil
}

func (me BasePluginLoader) Stop(logger comm.Logger) error {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	if !me.started {
		logger.Info().Msg("stopped, already")
		return nil
	}
	me.started = false

	errs := comm.NewErrorGroup(false)
	ns := me.Namespace()

	for _, plugin := range me.plugins {
		if err := StopPlugin(ns, plugin, logger); err != nil {
			errs.Add(err)
		}
	}

	if errs.HasError() {
		return errs
	}
	return nil
}
