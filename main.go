// get-headers prints the headers from GET requesting a URL
package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	netURL "net/url"
	"os"

	"github.com/carlmjohnson/get-headers/prettyprint"
)

const usage = `Usage of get-headers:

	get-headers <url>...
		Gets the URLs and prints their headers alphabetically.
		Repeated headers are printed with an asterisk.
	get-headers (-h|--help)
		Print this help message.
`

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(2)
	}
}

// Sentinal error to let us know if we're ignoring a redirect
var errRedirect = errors.New("redirected")

// Don't follow redirects
func checkRedirect(req *http.Request, via []*http.Request) error {
	return errRedirect
}

func main() {
	for _, url := range flag.Args() {
		client := http.Client{
			CheckRedirect: checkRedirect,
		}
		resp, err := client.Get(url)

		// Ignore the error if it's just our errRedirect
		switch urlErr, ok := err.(*netURL.Error); {
		case err == nil:
		case ok && urlErr.Err == errRedirect:
		default:
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		resp.Body.Close()
		fmt.Println("GET", url)
		fmt.Println(resp.Proto, resp.Status, "\n")
		fmt.Println(prettyprint.ResponseHeader(resp.Header))
	}
}
