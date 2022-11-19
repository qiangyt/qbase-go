package test

import (
	"testing"

	"github.com/fastgh/go-comm/v2"
	qplugin "github.com/qiangyt/qbase-go/plugin"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_ListExternalPlugins_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// /plugins/a01
	comm.WriteFileTextP(fs, "/plugins/a01/plugin.manifest.yml", `
kind: hosts
name: PluginA
version_major: 0
version_minor: 1
`)
	comm.WriteFileTextP(fs, "/plugins/a01/plugin.go", "package plugin")

	// /plugins/a02
	comm.WriteFileTextP(fs, "/plugins/a02/plugin.manifest.yml", `
kind: hosts
name: PluginA
version_major: 0
version_minor: 2
`)
	comm.WriteFileTextP(fs, "/plugins/a02/plugin.go", "package plugin")

	// /plugins/a10
	comm.WriteFileTextP(fs, "/plugins/a10/plugin.manifest.yml", `
kind: hosts
name: PluginA
version_major: 1
version_minor: 0
`)
	comm.WriteFileTextP(fs, "/plugins/a10/plugin.go", "package plugin")

	// /plugins/a30, but no plugin.go
	comm.WriteFileTextP(fs, "/plugins/a30/plugin.manifest.yml", `
kind: hosts
name: PluginA
version_major: 3
version_minor: 0
`)

	// /plugins/a40, but no manifest
	comm.WriteFileTextP(fs, "/plugins/a40/plugin.go", "package plugin")

	// /plugins/b
	comm.WriteFileTextP(fs, "/plugins/b/plugin.manifest.yml", `
kind: tool
name: PluginB
version_major: 2
version_minor: 3
`)
	comm.WriteFileTextP(fs, "/plugins/b/plugin.go", "package plugin")

	plugins := qplugin.ListExternalPlugins(comm.NewDiscardLogger(), fs, "/plugins")
	a.Len(plugins, 2)

	pa := plugins[0]
	a.Equal("hosts", pa.Kind())
	a.Equal("plugina", pa.Name())
	a.Equal("go", pa.Language())
	a.Equal("/plugins/a10", pa.Dir())

	paMajor, paMinor := pa.Version()
	a.Equal(1, paMajor)
	a.Equal(0, paMinor)

	a.Equal("/plugins/a10/plugin.go", pa.CodeFile())
	a.False(pa.IsStarted())

	pb := plugins[1]
	a.Equal("tool", pb.Kind())
	a.Equal("pluginb", pb.Name())
	a.Equal("go", pb.Language())
	a.Equal("/plugins/b", pb.Dir())

	pbMajor, pbMinor := pb.Version()
	a.Equal(2, pbMajor)
	a.Equal(3, pbMinor)

	a.Equal("/plugins/b/plugin.go", pb.CodeFile())
	a.False(pb.IsStarted())
}
