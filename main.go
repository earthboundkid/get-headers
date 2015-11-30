// get-headers prints the headers from GET requesting a URL
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
)

//respHeader is a type for pretty printing response headers
type respHeader http.Header

func (h respHeader) String() string {
	// Sort headers; get max header string length
	sortedHeaderKeys := make([]string, 0, len(h))
	for header := range h {
		sortedHeaderKeys = append(sortedHeaderKeys, header)
	}
	sort.Strings(sortedHeaderKeys)

	// Use a tabwriter to pretty print the output to a buffer
	var (
		buf bytes.Buffer
		tw  = tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	)
	for _, headerKey := range sortedHeaderKeys {
		for i, headerValue := range h[headerKey] {
			// Flag repeated values with an asterisk
			asterisk := ""
			if i > 0 {
				asterisk = " *"
			}
			fmt.Fprintf(tw, "%s%s\t%s\t\n", headerKey, asterisk, headerValue)
		}
	}
	tw.Flush()
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
		fmt.Println("GET", url)
		fmt.Println(resp.Proto, resp.Status, "\n")
		fmt.Println(respHeader(resp.Header))
	}
}
