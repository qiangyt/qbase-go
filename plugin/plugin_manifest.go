package qplugin

import (
	"strings"

	"github.com/fastgh/go-comm/v2"
	"github.com/spf13/afero"
)

type PluginManifestT struct {
	Kind         PluginKind `mapstructure:"kind" yaml:"kind"`
	Name         string     `mapstructure:"name" yaml:"name"`
	VersionMajor int        `mapstructure:"version_major" yaml:"version_major"`
	VersionMinor int        `mapstructure:"version_minor" yaml:"version_minor"`
}

type PluginManifest = *PluginManifestT

func PluginManifestWithMap(manifestMap map[string]any) PluginManifest {
	r := comm.DecodeWithMapP(manifestMap, &comm.ConfigConfig{
		ErrorUnused:          true,
		ErrorUnset:           false,
		ZeroFields:           false,
		WeaklyTypedInput:     true,
		Squash:               true,
		IgnoreUntaggedFields: true,
	}, &PluginManifestT{}, nil)

	r.Name = strings.ToLower(r.Name)

	return r
}

func PluginManifestWithJsonFile(fs afero.Fs, manifestJsonFile string) PluginManifest {
	manifestMap := comm.MapFromJsonFileP(fs, manifestJsonFile, false)
	return PluginManifestWithMap(manifestMap)
}

func PluginManifestWithYamlFile(fs afero.Fs, manifestYamlFile string) PluginManifest {
	manifestMap := comm.MapFromYamlFileP(fs, manifestYamlFile, false)
	return PluginManifestWithMap(manifestMap)
}
