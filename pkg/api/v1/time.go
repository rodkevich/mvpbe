package v1

import (
	"time"
)

// TimeNow for in application responses usage
var TimeNow = time.Now().Truncate(time.Microsecond)
