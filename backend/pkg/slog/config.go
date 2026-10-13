package slog

type Config struct {
	DisableColors   bool   `default:"true"`
	FullTimestamp   bool   `default:"true"`
	SetReportCaller bool   `default:"true"`
	Output          string `default:"stdout"`
	PathToFile      string `default:""`
	Level           string `default:"info"`
}
