package prettyprint

import "fmt"

// A petabyte is 2^50, so a float64, which holds 2^53 ints, is big
// enough for whatever purposes.

// Size is wrapper for humanizing byte sizes. Sizes are reported in
// base-2 equivalents, not base-10, i.e. 1 KB = 1024 bytes.
type Size float64

func (size Size) String() string {
	const (
		kilobyte = 1 << (10 * (iota + 1))
		megabyte
		gigabyte
		terabyte
		petabyte
	)

	format := "%.f B"

	switch {
	case size >= petabyte:
		format = "%3.1f PB"
		size /= terabyte
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
