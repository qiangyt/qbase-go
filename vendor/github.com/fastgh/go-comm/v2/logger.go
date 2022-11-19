package comm

import (
	"io"
	"os"
	"path/filepath"

	"github.com/fastgh/go-event"
	eventloggers "github.com/fastgh/go-event/loggers/phuslu"
	plog "github.com/phuslu/log"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerT struct {
	plog.Logger
	parent           Logger
	lumberjackLogger *lumberjack.Logger
}

type (
	Logger     = *LoggerT
	LogEntry   = *plog.Entry
	LogContext = LogEntry
)

var TraceId atomic.Int64

func NewLogContext(generateNewTraceId bool) LogContext {
	r := plog.NewContext(nil)
	if generateNewTraceId {
		r.Int64("traceId", TraceId.Add(1))
	}
	return r
}

// see lumberjack.Logger
type LoggerConfigT struct {
	MaxSize    int  `json:"max_size" yaml:"maxsize"`
	MaxAge     int  `json:"max_age" yaml:"maxage"`
	MaxBackups int  `json:"max_backups" yaml:"maxbackups"`
	LocalTime  bool `json:"local_time" yaml:"localtime"`
	Compress   bool `json:"compress" yaml:"compress"`
}

type LoggerConfig = *LoggerConfigT

func (me Logger) NewSubLogger(lctx LogContext) Logger {
	r := *me
	r.parent = me
	if lctx != nil {
		r.Context = lctx.Value()
	}
	return &r
}

func (me Logger) Parent() Logger {
	return me.parent
}

func (me Logger) Close() {
	if me.lumberjackLogger != nil {
		me.lumberjackLogger.Close()
	}
}

func (me Logger) Error(err any) LogEntry {
	r := me.Logger.Error()
	eventloggers.PhusluMarshalAnyError(r, err)
	return r
}

func NewLoggerP(console io.Writer, config LoggerConfig, fileName string) Logger {
	r, err := NewLogger(console, config, fileName)
	if err != nil {
		panic(err)
	}
	return r
}

// / verbose: log to console if true
func NewLogger(console io.Writer, config LoggerConfig, fileName string) (Logger, error) {
	logD := filepath.Dir(fileName)
	if err := os.MkdirAll(logD, os.ModePerm); err != nil {
		return nil, errors.Wrapf(err, "create directory: %s", logD)
	}

	// we use lumberjack instead of phuslu/log/FileWritter as said by phuslu/log that:
	// 	"FileWriter creates a symlink to the current logging file, it requires administrator privileges on Windows."
	//  'administrator privileges' is not acceptable for our scenario
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		LocalTime:  config.LocalTime,
		Compress:   config.Compress,
	}
	fileW := &plog.IOWriter{
		Writer: lumberjackLogger,
	}

	writers := plog.MultiEntryWriter{fileW}

	if console != nil {
		consoleW := &plog.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    false,
			EndWithMessage: true,
			Writer:         console,
		}

		writers = append(writers, consoleW)
	}

	return &LoggerT{
		Logger: plog.Logger{
			Level:  plog.InfoLevel,
			Caller: 3,
			Writer: &writers,
		},
		parent:           nil,
		lumberjackLogger: lumberjackLogger,
	}, nil
}

func NewDiscardLogger() Logger {
	return &LoggerT{
		Logger: plog.Logger{
			Level: plog.FatalLevel,
			Writer: plog.IOWriter{
				Writer: io.Discard, // log is off by default
			},
		},
		parent: nil,
	}
}

func IsDiscardLogger(logger Logger) bool {
	return logger.Logger.Writer.(plog.IOWriter).Writer == io.Discard
}

type EventLoggerT struct {
	target Logger
}

type EventLogger = *EventLoggerT

func NewEventLogger(target Logger) EventLogger {
	return &EventLoggerT{
		target: target,
	}
}

func (me EventLogger) Target() Logger { return me.target }

func (me EventLogger) LogDebug(enm event.LogEnum, hub string, topic string, lsner string) {
	me.target.Debug().Str("hub", hub).Str("topic", topic).Str("listener", lsner).Msg(enm.String())
}

func (me EventLogger) LogInfo(enm event.LogEnum, hub string, topic string, lsner string) {
	me.target.Info().Str("hub", hub).Str("topic", topic).Str("listener", lsner).Msg(enm.String())
}

func (me EventLogger) LogError(enm event.LogEnum, hub string, topic string, lsner string, err any) {
	entry := me.target.Error(err).Str("hub", hub).Str("topic", topic).Str("listener", lsner)
	entry.Msg(enm.String())
}

func (me EventLogger) LogEventDebug(enm event.LogEnum, lsner string, evnt event.Event) {
	me.target.Debug().Object("event", evnt).Str("listener", lsner).Msg(enm.String())
}

func (me EventLogger) LogEventInfo(enm event.LogEnum, lsner string, evnt event.Event) {
	me.target.Info().Object("event", evnt).Str("listener", lsner).Msg(enm.String())
}

func (me EventLogger) LogEventError(enm event.LogEnum, lsner string, evnt event.Event, err any) {
	entry := me.target.Error(err).Object("event", evnt).Str("listener", lsner)
	entry.Msg(enm.String())
}
