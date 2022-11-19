package qplugin

import (
	"path/filepath"

	"github.com/fastgh/go-comm/v2"
	"github.com/spf13/afero"
)

type FsPluginLoaderT struct {
	BasePluginLoaderT

	fs  afero.Fs
	dir string
}

type FsPluginLoader = *FsPluginLoaderT

func NewLocalPluginLoader(logger comm.Logger, fs afero.Fs, dir string) PluginLoader {
	return NewFsPluginLoader(logger, fs, dir, "local")
}

func NewRemotePluginLoader(logger comm.Logger, fs afero.Fs, dir string) PluginLoader {
	return NewFsPluginLoader(logger, fs, dir, "remote")
}

func NewFsPluginLoader(logger comm.Logger, fs afero.Fs, dir string, namespace string) PluginLoader {
	return &FsPluginLoaderT{
		BasePluginLoaderT: *NewPluginLoader(namespace),
		fs:                fs,
		dir:               filepath.Join(dir, namespace),
	}
}

func (me FsPluginLoader) Start(logger comm.Logger) error {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	if me.started {
		logger.Info().Msg("started, already")
		return nil
	}

	errs := comm.NewErrorGroup(false)
	ns := me.Namespace()

	for _, plugin := range ListExternalPlugins(logger, me.fs, filepath.Join(me.dir, me.namespace)) {
		me.Register(plugin)
	}

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
