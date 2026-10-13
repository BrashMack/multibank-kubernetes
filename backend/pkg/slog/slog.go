package slog

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger(cfg Config) {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: cfg.DisableColors,
		FullTimestamp: cfg.FullTimestamp,
	})

	switch cfg.Level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	switch cfg.Output {
	case "stdout":
		logrus.SetOutput(os.Stdout)
	case "stderr":
		logrus.SetOutput(os.Stderr)
	case "file":
		file, err := os.Open(cfg.PathToFile)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.SetOutput(file)
	}

	logrus.SetReportCaller(cfg.SetReportCaller)
}
