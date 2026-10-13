package jwt

import "time"

type Config struct {
	TTL    time.Duration `default:"168h"`
	Secret string        `required:"true"`
}
