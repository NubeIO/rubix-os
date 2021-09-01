package logger

import (
	log "github.com/sirupsen/logrus"
)

func SetLogger(logLevel string) {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	//log.SetReportCaller(true)
	log.AddHook(&MqttFieldHook{})
	log.SetLevel(getLogLevel(logLevel))
}

func getLogLevel(logLevel string) log.Level {
	var level log.Level
	switch logLevel {
	case "DEBUG":
		level = log.DebugLevel
	case "WARN":
		level = log.WarnLevel
	case "ERROR":
		level = log.ErrorLevel
	default:
		level = log.InfoLevel
	}
	return level
}

type MqttFieldHook struct {
}

func (h *MqttFieldHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *MqttFieldHook) Fire(e *log.Entry) error {
	//TODO: extend feature for MQTT streaming
	//fmt.Println(e.Level, e.Time, e.Message)
	return nil
}
