package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

func New() *logrus.Logger {
	// Send all logs to nowhere by default
	logrus.SetOutput(io.Discard)

	// Send logs with level higher than warning to stderr
	logrus.AddHook(&writer.Hook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	// Send info and debug logs to stdout
	logrus.AddHook(&writer.Hook{
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})

	logrus.SetFormatter(new(logrus.JSONFormatter))
	return logrus.New()
}
