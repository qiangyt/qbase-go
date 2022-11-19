package qplugin

import (
	"path/filepath"
	"sync"

	"github.com/fastgh/go-comm/v2"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type ExternalPluginContext interface {
	Init(logger comm.Logger, fs afero.Fs, codeFile string)
	Start() any
	Stop() any
}

type ExternalPluginT struct {
	kind     PluginKind
	name     string
	language string
	dir      string

	versionMajor int
	versionMinor int
	codeFile     string

	started bool

	context ExternalPluginContext

	mutex sync.RWMutex
}

type ExternalPlugin = *ExternalPluginT

func (me ExternalPlugin) Name() string {
	return me.name
}

func (me ExternalPlugin) IsStarted() bool {
	me.mutex.RLock()
	defer me.mutex.RUnlock()

	return me.started
}

func (me ExternalPlugin) Start(logger comm.Logger) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	if me.started {
		return
	}

	me.context.Start()

	me.started = true
}

func (me ExternalPlugin) Kind() PluginKind {
	return me.kind
}

func (me ExternalPlugin) Stop(logger comm.Logger) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	if !me.started {
		return
	}

	me.context.Stop()

	me.started = false
}

func (me ExternalPlugin) Version() (major int, minor int) {
	return me.versionMajor, me.versionMinor
}

func (me ExternalPlugin) Language() string {
	return me.language
}

func (me ExternalPlugin) Dir() string {
	return me.dir
}

func (me ExternalPlugin) CodeFile() string {
	return me.codeFile
}

func ResolveExternalPlugin(logger comm.Logger, fs afero.Fs, pluginDir string) (result ExternalPlugin) {
	defer func() {
		if p := recover(); p != nil {
			logger.Error(p).Str("pluginDir", pluginDir).Msg("failed to resolve external plugin")
			result = nil
		}
	}()

	var mf PluginManifest

	yamlF := filepath.Join(pluginDir, "plugin.manifest.yml")
	if comm.FileExistsP(fs, yamlF) {
		mf = PluginManifestWithYamlFile(fs, yamlF)
	} else {
		yamlF = filepath.Join(pluginDir, "plugin.manifest.yaml")
		if comm.FileExistsP(fs, yamlF) {
			mf = PluginManifestWithYamlFile(fs, yamlF)
		} else {
			jsonF := filepath.Join(pluginDir, "plugin.manifest.json")
			if comm.FileExistsP(fs, jsonF) {
				mf = PluginManifestWithJsonFile(fs, jsonF)
			}
		}
	}

	if mf == nil {
		return nil
	}

	language := PLUGIN_LANG_GO
	context := NewExternalGoPluginContext()
	codeFile := filepath.Join(pluginDir, "plugin.go")
	if exists, err := comm.FileExists(fs, codeFile); err != nil || !exists {
		return nil
	}

	context.Init(logger, fs, codeFile)

	result = &ExternalPluginT{
		kind:         mf.Kind,
		name:         mf.Name,
		language:     language,
		dir:          pluginDir,
		versionMajor: mf.VersionMajor,
		versionMinor: mf.VersionMinor,
		codeFile:     codeFile,
		started:      false,
		context:      context,
		mutex:        sync.RWMutex{},
	}

	return
}

func ListExternalPlugins(logger comm.Logger, afs afero.Fs, baseDir string) []ExternalPlugin {
	pluginDirOrFiles, err := afero.ReadDir(afs, baseDir)
	if err != nil {
		panic(errors.Wrapf(err, "read plugins directories: %s", baseDir))
	}

	r := comm.NewOrderedMap[ExternalPlugin](nil)

	for _, dirOrFile := range pluginDirOrFiles {
		if !dirOrFile.IsDir() {
			continue
		}

		fName := dirOrFile.Name()
		fBase := filepath.Base(fName)
		if fBase == "." || fBase == ".." {
			continue
		}

		pluginDir := filepath.Join(baseDir, fName)
		p := ResolveExternalPlugin(logger, afs, pluginDir)
		if p == nil {
			continue
		}

		name := p.Name()
		existing := r.Get(name)
		if existing == nil {
			r.Put(name, p)
		} else if p.versionMajor > existing.versionMajor {
			r.Put(name, p)
		} else if p.versionMajor == existing.versionMajor && p.versionMinor > existing.versionMinor {
			r.Put(name, p)
		}
	}

	return r.Values()
}
