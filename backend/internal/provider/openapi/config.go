package openapi

import (
	"time"

	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

type Retry struct {
	Codes    []int
	Period   time.Duration `default:"30s"`
	Deadline time.Duration `default:"5m"`
}

type ClientConfig struct {
	Retry               Retry
	ClientCert          string
	ClientKey           string
	RootCACert          string
	Insecure            bool          `default:"true"`
	MaxIdleConns        int           `default:"100"`
	MaxIdleConnsPerHost int           `default:"10"`
	IdleConnTimeout     time.Duration `default:"90s"`
	Timeout             time.Duration `default:"10s"`
	HiddenHeadersInLog  []string
}

type Config struct {
	Disabled     bool
	BaseUrls     map[enums.Bank]string
	TeamID       string
	TeamSecretID string
	Client       ClientConfig
	BankToCode   map[enums.Bank]string
}
