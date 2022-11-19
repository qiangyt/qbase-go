package event

type LogLevel int8

const (
	LogLevelAnyway LogLevel = iota

	LogLevelDebug
	LogLevelInfo
	LogLevelError

	LogLevelSilient LogLevel = 127
)

type Logger interface {
	LogDebug(enm LogEnum, hub string, topic string, lsner string)
	LogInfo(enm LogEnum, hub string, topic string, lsner string)
	LogError(enm LogEnum, hub string, topic string, lsner string, err any)

	LogEventDebug(enm LogEnum, lsner string, evnt Event)
	LogEventInfo(enm LogEnum, lsner string, evnt Event)
	LogEventError(enm LogEnum, lsner string, evnt Event, err any)
}
