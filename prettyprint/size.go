package prettyprint

import "fmt"

// Size is wrapper for humanizing byte sizes
type Size int64

func (size Size) String() string {
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
