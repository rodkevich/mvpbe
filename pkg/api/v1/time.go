package v1

import (
	"time"
)

// TimeNow for in application responses
func TimeNow() time.Time { return time.Now().Truncate(time.Microsecond) }
