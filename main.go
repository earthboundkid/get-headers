// get-headers prints the headers from GET requesting a URL
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
)

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

	// Make format strings as long as longest header
	fstrOuter := " %-" + strconv.Itoa(max+1) + "s"
	fstrInner := "%" + strconv.Itoa(max+2) + "s"

	buf := &bytes.Buffer{}
	for _, headerKey := range sortedHeaderKeys {
		fmt.Fprintf(buf, fstrOuter, headerKey)
		for i := range h[headerKey] {
			if i > 0 {
				fmt.Fprintf(buf, fstrInner, "")
			}
			fmt.Fprintf(buf, "%s\n", h[headerKey][i])
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
