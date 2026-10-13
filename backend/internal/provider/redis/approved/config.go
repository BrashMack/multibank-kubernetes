package approved

import "time"

type ClientConfig struct {
	Host            string
	Port            string
	Password        string
	DB              int
	UseTLS          bool
	DialTimeout     time.Duration `default:"5s"`
	ReadTimeout     time.Duration `default:"3s"`
	WriteTimeout    time.Duration `default:"3s"`
	PoolSize        int           `default:"10"`
	MinIdleConns    int           `default:"2"`
	ConnMaxIdleTime time.Duration `default:"30m"`
	ConnMaxLifetime time.Duration `default:"2h"`
}

type Config struct {
	Client            ClientConfig
	DefaultExpiration time.Duration `default:"2160h"`
}
