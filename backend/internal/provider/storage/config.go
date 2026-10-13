package storage

import "time"

type Config struct {
	Host                  string
	Database              string
	User                  string
	Password              string
	Port                  string
	DSN                   string
	MaxConnLifetime       time.Duration `default:"4h"`
	MaxConnLifetimeJitter time.Duration `default:"10m"`
	MaxConnIdleTime       time.Duration `default:"1h"`
	MaxConns              int32         `default:"100"`
	MinConns              int32         `default:"10"`
	HealthCheckPeriod     time.Duration `default:"1m"`
}
