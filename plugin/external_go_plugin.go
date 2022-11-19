package qplugin

import (
	"reflect"

	"github.com/fastgh/go-comm/v2"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type ExternalGoPluginContextT struct {
	interpreter *interp.Interpreter

	startFunc *reflect.Value
	stopFunc  *reflect.Value
}

type ExternalGoPluginContext = *ExternalGoPluginContextT

func NewExternalGoPluginContext() ExternalGoPluginContext {
	return &ExternalGoPluginContextT{
		interpreter: nil,
		startFunc:   nil,
		stopFunc:    nil,
	}
}

func resolveExternalGoPluginFunc(logger comm.Logger, interpreter *interp.Interpreter, funcName string) *reflect.Value {
	r, err := interpreter.Eval(funcName)
	if err != nil {
		logger.Error(err).Msg("failed to eval " + funcName)
		return nil
	}
	if comm.IsPrimitiveReflectValue(r) {
		logger.Error(err).Msg(funcName + " is a primitive value instead of a function")
		return nil
	}
	if r.IsNil() {
		logger.Error(err).Msg("symbol not found: " + funcName)
		return nil
	}
	if r.Kind() != reflect.Func {
		logger.Error(err).Msg(funcName + " is not a function")
		return nil
	}

	return &r
}

func (me ExternalGoPluginContext) Init(logger comm.Logger, fs afero.Fs, codeFile string) {
	logCtx := comm.NewLogContext(false)
	logCtx.Str("codeFile", codeFile)
	logger = logger.NewSubLogger(logCtx)

	me.interpreter = interp.New(interp.Options{})
	if err := me.interpreter.Use(stdlib.Symbols); err != nil {
		panic(errors.Wrapf(err, "use stdlib failed: %s", codeFile))
	}

	code := comm.ReadFileTextP(fs, codeFile)
	_, err := me.interpreter.Eval(code)
	if err != nil {
		panic(errors.Wrapf(err, "eval %s", codeFile))
	}

	me.startFunc = resolveExternalGoPluginFunc(logger, me.interpreter, "plugin.PluginStart")
	me.stopFunc = resolveExternalGoPluginFunc(logger, me.interpreter, "plugin.PluginStop")
}

func (me ExternalGoPluginContext) GetStartFunc() *reflect.Value {
	return me.startFunc
}

func (me ExternalGoPluginContext) Start() any {
	if me.startFunc == nil {
		return ""
	}
	return me.startFunc.Call([]reflect.Value{})
}

func (me ExternalGoPluginContext) GetStopFunc() *reflect.Value {
	return me.stopFunc
}

func (me ExternalGoPluginContext) Stop() any {
	if me.stopFunc == nil {
		return ""
	}
	return me.stopFunc.Call([]reflect.Value{})
}
