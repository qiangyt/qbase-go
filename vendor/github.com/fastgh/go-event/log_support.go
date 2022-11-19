package event

type LevelFilteringLoggerT struct {
	target Logger
	Level  LogLevel
}

type LevelFilteringLogger = *LevelFilteringLoggerT

func NewLevelFilteringLogger(level LogLevel, target Logger) LevelFilteringLogger {
	return &LevelFilteringLoggerT{
		Level:  level,
		target: target,
	}
}

func (me LevelFilteringLogger) Target() Logger { return me.target }

func (me LevelFilteringLogger) LogDebug(enm LogEnum, hub string, topic string, lsner string) {
	if me.Level <= LogLevelDebug {
		me.target.LogDebug(enm, hub, topic, lsner)
	}
}

func (me LevelFilteringLogger) LogInfo(enm LogEnum, hub string, topic string, lsner string) {
	if me.Level <= LogLevelInfo {
		me.target.LogInfo(enm, hub, topic, lsner)
	}
}

func (me LevelFilteringLogger) LogError(enm LogEnum, hub string, topic string, lsner string, err any) {
	if me.Level <= LogLevelError {
		me.target.LogError(enm, hub, topic, lsner, err)
	}
}

func (me LevelFilteringLogger) LogEventDebug(enm LogEnum, lsner string, evnt Event) {
	if me.Level <= LogLevelDebug {
		me.target.LogEventDebug(enm, lsner, evnt)
	}
}

func (me LevelFilteringLogger) LogEventInfo(enm LogEnum, lsner string, evnt Event) {
	if me.Level <= LogLevelInfo {
		me.target.LogEventInfo(enm, lsner, evnt)
	}
}

func (me LevelFilteringLogger) LogEventError(enm LogEnum, lsner string, evnt Event, err any) {
	if me.Level <= LogLevelError {
		me.target.LogEventError(enm, lsner, evnt, err)
	}
}
