package test

import (
	"reflect"
	"testing"

	"github.com/fastgh/go-comm/v2"
	qplugin "github.com/qiangyt/qbase-go/plugin"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_ExternalGoPlugin_compileError(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.WriteFileTextP(fs, "/plugin.go", "nothing")

	p := qplugin.NewExternalGoPluginContext()
	a.Panics(func() {
		p.Init(comm.NewDiscardLogger(), fs, "/plugin.go")
	})
}

func Test_ExternalGoPlugin_full(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.WriteFileTextP(fs, "/plugin.go", `
	package plugin

	func PluginStart() string {
		return "start"
	}

	func PluginStop() string {
		return "stop"
	}
	`)

	p := qplugin.NewExternalGoPluginContext()
	p.Init(comm.NewDiscardLogger(), fs, "/plugin.go")

	a.NotNil(p.GetStartFunc())
	a.Equal("start", p.Start().([]reflect.Value)[0].String())

	a.NotNil(p.GetStopFunc())
	a.Equal("stop", p.Stop().([]reflect.Value)[0].String())
}

func Test_ExternalGoPlugin_noStart(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.WriteFileTextP(fs, "/plugin.go", `
	package plugin

	func PluginStop() string {
		return "stop"
	}
	`)

	p := qplugin.NewExternalGoPluginContext()
	p.Init(comm.NewDiscardLogger(), fs, "/plugin.go")

	a.Nil(p.GetStartFunc())
	a.Empty(p.Start())

	a.NotNil(p.GetStopFunc())
	a.Equal("stop", p.Stop().([]reflect.Value)[0].String())
}

func Test_ExternalGoPlugin_noStop(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.WriteFileTextP(fs, "/plugin.go", `
	package plugin

	func PluginStart() string {
		return "start"
	}
	`)

	p := qplugin.NewExternalGoPluginContext()
	p.Init(comm.NewDiscardLogger(), fs, "/plugin.go")

	a.NotNil(p.GetStartFunc())
	a.Equal("start", p.Start().([]reflect.Value)[0].String())

	a.Nil(p.GetStopFunc())
	a.Empty(p.Stop())
}

func Test_ExternalGoPlugin_startIsNotFunction(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.WriteFileTextP(fs, "/plugin.go", `
	package plugin

	const PluginStart = ""
	`)

	p := qplugin.NewExternalGoPluginContext()
	p.Init(comm.NewDiscardLogger(), fs, "/plugin.go")

	a.Nil(p.GetStartFunc())
	a.Empty(p.Start())
}

func Test_ExternalGoPlugin_stopIsNotFunction(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.WriteFileTextP(fs, "/plugin.go", `
	package plugin

	const PluginStop = ""
	`)

	p := qplugin.NewExternalGoPluginContext()
	p.Init(comm.NewDiscardLogger(), fs, "/plugin.go")

	a.Nil(p.GetStopFunc())
	a.Empty(p.Stop())
}
