package prettyprint

import (
	"fmt"
	"time"
)

// Duration is a wrapper for humanizing time.Duration
type Duration time.Duration

func (duration Duration) String() string {
	d := time.Duration(duration)
	minutes := d / time.Minute
	seconds := d % time.Minute / time.Second
	milli := d % time.Second / time.Millisecond
	micro := d % time.Millisecond / time.Microsecond

	switch {
	case minutes > 0:
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	case seconds > 0:
		return fmt.Sprintf("%ds %dms", seconds, milli)
	case milli > 0:
		return fmt.Sprintf("%dms %dµs", milli, micro)
	}

	return d.String()
}
