package gcloud

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
)

type GCloudFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
}

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

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = logrus.DefaultTimestampFormat
	}

	data["time"] = entry.Time.Format(timestampFormat)
	data["message"] = entry.Message
	data["severity"] = entry.Level.String()

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Logger failed to marshal fields to JSON, %v. Message: %v", err, entry.Message)
	}
	return append(serialized, '\n'), nil
}

func prefixFieldClashes(data logrus.Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["message"]; ok {
		data["fields.message"] = m
	}

	if l, ok := data["severity"]; ok {
		data["fields.severity"] = l
	}
}
