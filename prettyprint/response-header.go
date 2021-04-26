package prettyprint

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"text/tabwriter"
)

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
	for _, headerKey := range sortedHeaderKeys {
		for i, headerValue := range h[headerKey] {
			// Flag repeated values with an asterisk
			asterisk := ""
			if i > 0 {
				asterisk = " *"
			}

			// Prevent long lines by breaking at "; "
			if len(headerValue) > 50 {
				headerValue = strings.Replace(headerValue, "; ", ";\n...\t", -1)
			}
			fmt.Fprintf(tw, "%s%s\t%s\n", headerKey, asterisk, headerValue)
		}
	}
	tw.Flush()
	return buf.String()
}
