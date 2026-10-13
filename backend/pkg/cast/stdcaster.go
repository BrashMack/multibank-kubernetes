package cast

import (
	"time"
)

func init() {
	MustRegisterCaster(func(src string, opt Options) (time.Time, error) {
		return WrapCastError(time.Parse(opt.TimeLayout, src))
	})

	MustRegisterCaster(func(src string, _ Options) (time.Duration, error) {
		return WrapCastError(time.ParseDuration(src))
	})
}
