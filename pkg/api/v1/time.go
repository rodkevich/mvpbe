package v1

import (
	"time"
)

var (
	// TimeNow for in application responses usage
	TimeNow = time.Now().Truncate(time.Microsecond)
)
