// get-headers prints the headers from GET requesting a URL
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sort"
)

//respHeader is a type for pretty printing response headers
type respHeader http.Header

func (h respHeader) String() string {
	// Sort headers; get max header string length
	sortedHeaderKeys := make([]string, 0, len(h))
	max := 0
	for header := range h {
		sortedHeaderKeys = append(sortedHeaderKeys, header)
		if len(header) > max {
			max = len(header)
		}
	}
	sort.Strings(sortedHeaderKeys)

	buf := &bytes.Buffer{}
	for _, headerKey := range sortedHeaderKeys {
		fmt.Fprintf(buf, " %-*s", max+1, headerKey)
		headerValues := h[headerKey]
		for i := range headerValues {
			if i > 0 {
				fmt.Fprintf(buf, "%*s", max+2, "")
			}
			fmt.Fprintf(buf, "%s\n", headerValues[i])
		}
	}
	return buf.String()
}

func main() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		resp.Body.Close()
		fmt.Print(respHeader(resp.Header))
	}
}
