package prettyprint

import "fmt"

// Size is wrapper for humanizing byte sizes
type Size float64

func (size Size) String() string {
	const (
		kilobyte = 1 << (10 * (iota + 1))
		megabyte
		gigabyte
		terabyte
	)

	format := "%.f"

	switch {
	case size >= terabyte:
		format = "%3.1f TB"
		size /= terabyte
	case size >= gigabyte:
		format = "%3.1f GB"
		size /= gigabyte
	case size >= megabyte:
		format = "%3.1f MB"
		size /= megabyte
	case size >= kilobyte:
		format = "%3.1f KB"
		size /= kilobyte
	}
	return fmt.Sprintf(format, size)
}
