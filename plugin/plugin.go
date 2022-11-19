package qplugin

import (
	"github.com/fastgh/go-comm/v2"
)

type PluginLang = string

const (
	PLUGIN_LANG_GO         = "go"
	PLUGIN_LANG_JAVASCRIPT = "javascript"
	PLUGIN_LANG_SHELL      = "shell"
)

type PluginKind = string

type Plugin interface {
	Name() string
	Kind() PluginKind
	Start(logger comm.Logger)
	Stop(logger comm.Logger)
	Version() (major int, minor int)
}

type PluginLoader interface {
	Namespace() string
	Plugins() map[string]Plugin
	Start(logger comm.Logger) error
	Stop(logger comm.Logger) error
}
