package prettyprint

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"golang.org/x/term"
)

var isTTY = term.IsTerminal(int(os.Stdout.Fd()))

// ResponseHeader is a type for pretty printing net/http response headers
type ResponseHeader http.Header

// String formats headers by outputting headers and values in equal columns, sorted alphabetically by header. Repeated headers are marked with an asterisk.
func (h ResponseHeader) String() string {
	// Sort headers; get max header string length
	sortedHeaderKeys := make([]string, 0, len(h))
	for header := range h {
		sortedHeaderKeys = append(sortedHeaderKeys, header)
	}
	sort.Strings(sortedHeaderKeys)

	// Use a tabwriter to pretty print the output to a buffer
	var (
		buf strings.Builder
		tw  = tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	)
	var bold, reset, faint string
	wrapline := ";\n...\t"
	if isTTY {
		reset = "\x1b[0m"
		bold = "\x1b[1m"
		faint = "\x1b[2m"
		wrapline = ";\n" + faint + "..." + reset + "\t"
	}
	for _, headerKey := range sortedHeaderKeys {
		for i, headerValue := range h[headerKey] {
			// Flag repeated values with an asterisk
			asterisk := ""
			if i > 0 {
				asterisk = " *"
			}

			// Prevent long lines by breaking at "; "
			if len(headerValue) > 50 {
				headerValue = strings.ReplaceAll(headerValue, "; ", wrapline)
			}
			fmt.Fprintf(tw, "%s%s%s%s\t%s\n",
				bold, headerKey, reset, asterisk, headerValue)
		}
	}
	tw.Flush()
	return buf.String()
}
