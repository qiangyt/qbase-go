package event

type ListenerLoggerT struct {
	logr  TopicLogger
	lsner string
}

type ListenerLogger = *ListenerLoggerT

func NewListenerLogger(lsner string, logr TopicLogger) ListenerLogger {
	return &ListenerLoggerT{
		logr:  logr,
		lsner: lsner,
	}
}

func (me ListenerLogger) Listener() string { return me.lsner }

func (me ListenerLogger) LogDebug(enm LogEnum) {
	me.logr.LogDebug(enm, me.lsner)
}

func (me ListenerLogger) LogInfo(enm LogEnum) {
	me.logr.LogInfo(enm, me.lsner)
}

func (me ListenerLogger) LogError(enm LogEnum, err any) {
	me.logr.LogError(enm, me.lsner, err)
}

func (me ListenerLogger) LogEventInfo(enm LogEnum, evnt Event) {
	me.logr.LogEventInfo(enm, me.lsner, evnt)
}

func (me ListenerLogger) LogEventDebug(enm LogEnum, evnt Event) {
	me.logr.LogEventDebug(enm, me.lsner, evnt)
}

func (me ListenerLogger) LogEventError(enm LogEnum, evnt Event, err any) {
	me.logr.LogEventError(enm, me.lsner, evnt, err)
}
