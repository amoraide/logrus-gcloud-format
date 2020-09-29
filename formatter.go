package gcloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	levelMap = map[logrus.Level]string{
		logrus.TraceLevel: "DEBUG",
		logrus.DebugLevel: "DEBUG",
		logrus.InfoLevel:  "INFO",
		logrus.WarnLevel:  "WARNING",
		logrus.ErrorLevel: "ERROR",
		logrus.FatalLevel: "FATAL",
		logrus.PanicLevel: "ALERT",
	}

	fieldClashes = []string{
		"time",
		"message",
		"severity",
	}
)

// GCloudFormatter is a custom format for GCE
type GCloudFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
}

// Format function parses the logrus.Entry for relevant information to be formatted for GCE
func (f *GCloudFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	prefixFieldClashes(data)

	format := time.RFC3339Nano
	if f.TimestampFormat != "" {
		format = f.TimestampFormat
	}

	data["time"] = entry.Time.Format(format)
	data["message"] = entry.Message
	data["severity"] = levelMap[entry.Level]

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Logger failed to marshal fields to JSON, %v. Message: %v", err, entry.Message)
	}
	return append(serialized, '\n'), nil
}

func prefixFieldClashes(data logrus.Fields) {
	for _, fc := range fieldClashes {
		if d, ok := data[fc]; ok {
			data[fmt.Sprintf("fields.%s", fc)] = d
		}
	}
}
