package main

import (
	"fmt"
	"time"
)

func humanizeByteSize(size int64) string {
	const (
		_        = iota
		kilobyte = 1 << (10 * iota)
		megabyte
		gigabyte
		terabyte
	)

	format := "%.f"
	value := float32(size)

	switch {
	case size >= terabyte:
		format = "%3.1f TB"
		value = value / terabyte
	case size >= gigabyte:
		format = "%3.1f GB"
		value = value / gigabyte
	case size >= megabyte:
		format = "%3.1f MB"
		value = value / megabyte
	case size >= kilobyte:
		format = "%3.1f KB"
		value = value / kilobyte
	}
	return fmt.Sprintf(format, value)
}

func humanizeDuration(d time.Duration) string {
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
	case micro > 0:
		return d.String()
	}
	// I am assuming nanosecond timing is impossible
	return ""
}
