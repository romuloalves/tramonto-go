package files

import "time"

var now = func() time.Time {
	return time.Now()
}
