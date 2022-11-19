package event

type TopicLoggerT struct {
	topic string
	logr  HubLogger
}

type TopicLogger = *TopicLoggerT

func NewTopicLogger(topic string, logr HubLogger) TopicLogger {
	return &TopicLoggerT{
		logr:  logr,
		topic: topic,
	}
}

func (me TopicLogger) Topic() string { return me.topic }

func (me TopicLogger) LogDebug(enm LogEnum, lsner string) {
	me.logr.LogDebug(enm, me.topic, lsner)
}

func (me TopicLogger) LogInfo(enm LogEnum, lsner string) {
	me.logr.LogInfo(enm, me.topic, lsner)
}

func (me TopicLogger) LogError(enm LogEnum, lsner string, err any) {
	me.logr.LogError(enm, me.topic, lsner, err)
}

func (me TopicLogger) LogEventDebug(enm LogEnum, lsner string, evnt Event) {
	me.logr.LogEventDebug(enm, lsner, evnt)
}

func (me TopicLogger) LogEventInfo(enm LogEnum, lsner string, evnt Event) {
	me.logr.LogEventInfo(enm, lsner, evnt)
}

func (me TopicLogger) LogEventError(enm LogEnum, lsner string, evnt Event, err any) {
	me.logr.LogEventError(enm, lsner, evnt, err)
}
