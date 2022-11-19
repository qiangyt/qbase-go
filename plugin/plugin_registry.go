package qplugin

import (
	"fmt"
	"sync"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/fastgh/go-comm/v2"
)

type PluginRegistryT struct {
	loaders               map[string]PluginLoader
	plugins               []Plugin
	pluginsByKind         map[PluginKind]map[string]Plugin
	supportedKinds        hashset.Set
	supportedMajorVersion int
	mutex                 sync.RWMutex
}

type PluginRegistry = *PluginRegistryT

func NewPluginRegistry(supportedMajorVersion int, supportedKinds ...PluginKind) PluginRegistry {
	r := &PluginRegistryT{
		loaders:               map[string]PluginLoader{},
		plugins:               []Plugin{},
		pluginsByKind:         map[PluginKind]map[string]Plugin{},
		supportedKinds:        *comm.Slice2Set(supportedKinds...),
		supportedMajorVersion: supportedMajorVersion,
		mutex:                 sync.RWMutex{},
	}
	return r
}

func (me PluginRegistry) IsSupportedPluginKind(kind PluginKind) bool {
	return me.supportedKinds.Contains(kind)
}

func (me PluginRegistry) SupportedMajorVersion() int {
	return me.supportedMajorVersion
}

func (me PluginRegistry) ValidatePlugin(namespace string, plugin Plugin) error {
	name := plugin.Name()

	major, _ := plugin.Version()
	if major != me.supportedMajorVersion {
		return fmt.Errorf("expect plugin %s/%s major version is %d, but it is %d",
			namespace, name, me.supportedMajorVersion, major)
	}

	kind := plugin.Kind()
	if !me.IsSupportedPluginKind(kind) {
		return fmt.Errorf("plugin %s/%s claims unsupported plugin kind %s", namespace, name, kind)
	}

	return nil
}

func (me PluginRegistry) ByKind(kind PluginKind) map[string]Plugin {
	me.mutex.RLock()
	defer me.mutex.RUnlock()

	return me.pluginsByKind[kind]
}

func (me PluginRegistry) Init(logger comm.Logger) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	for ns, loader := range me.loaders {
		logCtx := comm.NewLogContext(false)
		logCtx.Str("namespace", ns)
		subLogger := logger.NewSubLogger(logCtx)

		subLogger.Info().Msg("starting plugin loader")
		err := loader.Start(logger)
		if err != nil {
			subLogger.Error(err).Msg("failed to start plugin loader")
		} else {
			subLogger.Info().Msg("started plugin loader")
		}
	}
}

func (me PluginRegistry) Destroy(logger comm.Logger) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	for ns, loader := range me.loaders {
		logCtx := comm.NewLogContext(false)
		logCtx.Str("namespace", ns)
		subLogger := logger.NewSubLogger(logCtx)

		subLogger.Info().Msg("stopping plugin loader")
		err := loader.Stop(logger)
		if err != nil {
			subLogger.Error(err).Msg("failed to stop plugin loader")
		} else {
			subLogger.Info().Msg("stopped plugin loader")
		}
	}

	me.loaders = map[string]PluginLoader{}
	me.plugins = []Plugin{}
	me.pluginsByKind = map[PluginKind]map[string]Plugin{}
}

func (me PluginRegistry) HasNamespace(ns string) bool {
	_, r := me.loaders[ns]
	return r
}

func (me PluginRegistry) Register(loader PluginLoader) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	ns := loader.Namespace()
	if len(ns) == 0 {
		panic(fmt.Errorf("namespace not specified: %+v", loader))
	}

	if existingLoader, alreadyRegistered := me.loaders[ns]; alreadyRegistered {
		panic(fmt.Errorf("plugin namespace %s is already registered by: %+v", ns, existingLoader))
	}

	newPlugins := loader.Plugins()

	pluginsByKind := comm.DeepCopyMap(me.pluginsByKind)

	for name, plugin := range newPlugins {
		if err := me.ValidatePlugin(ns, plugin); err != nil {
			panic(err)
		}

		kind := plugin.Kind()
		id := PluginId(ns, name)

		pluginsWithKind, found := pluginsByKind[kind]
		if !found {
			pluginsWithKind = map[string]Plugin{}
			pluginsByKind[kind] = pluginsWithKind
		}

		if existingPlugin, found := pluginsWithKind[name]; found {
			panic(fmt.Errorf("plugin %s has duplicated kind %s with %+v", id, kind, existingPlugin))
		}
		pluginsWithKind[name] = plugin
	}

	allPlugins := make([]Plugin, len(me.plugins), len(me.plugins)+len(newPlugins))
	allPlugins = append(allPlugins, me.plugins...)

	for _, plugin := range newPlugins {
		allPlugins = append(allPlugins, plugin)
	}

	me.plugins = allPlugins
	me.pluginsByKind = pluginsByKind
}
