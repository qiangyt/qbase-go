package qplugin

import (
	"fmt"

	"github.com/fastgh/go-comm/v2"
	"github.com/pkg/errors"
)

func PluginId(namespace string, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

func StartPlugin(namespace string, plugin Plugin, logger comm.Logger) (err error) {
	major, minor := plugin.Version()
	ver := fmt.Sprintf("%d/%d", major, minor)
	pluginId := PluginId(namespace, plugin.Name())

	defer func() {
		if p := recover(); p != nil {
			var err2 error
			var isErr bool
			if err2, isErr = p.(error); isErr {
				err = errors.Wrapf(err2, "start plugin: %s (version %s)", pluginId, ver)
			} else {
				err = fmt.Errorf("start plugin: %s (version %s), cause: %+v", pluginId, ver, p)
			}
		}
	}()

	logCtx := comm.NewLogContext(false)
	logCtx.Str("pluginId", pluginId).Str("version", ver)
	subLogger := logger.NewSubLogger(logCtx)

	subLogger.Info().Msg("starting")
	plugin.Start(logger)
	subLogger.Info().Msg("started")

	return
}

func StopPlugin(namespace string, plugin Plugin, logger comm.Logger) (err error) {
	major, minor := plugin.Version()
	ver := fmt.Sprintf("%d/%d", major, minor)
	pluginId := PluginId(namespace, plugin.Name())

	defer func() {
		if p := recover(); p != nil {
			var err2 error
			var isErr bool
			if err2, isErr = p.(error); isErr {
				err = errors.Wrapf(err2, "stop plugin: %s (version %s)", pluginId, ver)
			} else {
				err = fmt.Errorf("stop plugin: %s (version %s), cause: %+v", pluginId, ver, p)
			}
		}
	}()

	logCtx := comm.NewLogContext(false)
	logCtx.Str("pluginId", pluginId).Str("version", ver)
	subLogger := logger.NewSubLogger(logCtx)

	subLogger.Info().Msg("stopping")
	plugin.Stop(logger)
	subLogger.Info().Msg("stopped")

	return
}
