package event

type HubLoggerT struct {
	hub  string
	logr Logger
}

type HubLogger = *HubLoggerT

func NewHubLogger(hub string, logr Logger) HubLogger {
	return &HubLoggerT{
		hub:  hub,
		logr: logr,
	}
}

func (me HubLogger) Hub() string {
	return me.hub
}

func (me HubLogger) LogDebug(enm LogEnum, topic string, lsner string) {
	if me.logr != nil {
		me.logr.LogDebug(enm, me.hub, topic, lsner)
	}
}

func (me HubLogger) LogInfo(enm LogEnum, topic string, lsner string) {
	if me.logr != nil {
		me.logr.LogInfo(enm, me.hub, topic, lsner)
	}
}

func (me HubLogger) LogError(enm LogEnum, topic string, lsner string, err any) {
	if me.logr != nil {
		me.logr.LogError(enm, me.hub, topic, lsner, err)
	}
}

func (me HubLogger) LogEventDebug(enm LogEnum, lsner string, evnt Event) {
	if me.logr != nil {
		me.logr.LogEventDebug(enm, lsner, evnt)
	}
}

func (me HubLogger) LogEventInfo(enm LogEnum, lsner string, evnt Event) {
	if me.logr != nil {
		me.logr.LogEventInfo(enm, lsner, evnt)
	}
}

func (me HubLogger) LogEventError(enm LogEnum, lsner string, evnt Event, err any) {
	me.logr.LogEventError(enm, lsner, evnt, err)
}
