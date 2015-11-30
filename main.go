// get-headers prints the headers from GET requesting a URL
package main

import (
	"errors"
	"fmt"
	"net/http"
	netURL "net/url"
	"os"

	"github.com/carlmjohnson/get-headers/prettyprint"
)

// Sentinal error to let us know if we're ignoring a redirect
var redirectError = errors.New("redirected")

// Don't follow redirects
func checkRedirect(req *http.Request, via []*http.Request) error {
	return redirectError
}

func main() {
	for _, url := range os.Args[1:] {
		client := http.Client{
			CheckRedirect: checkRedirect,
		}
		resp, err := client.Get(url)

		// Ignore the error if it's just our redirectError
		switch urlErr, ok := err.(*netURL.Error); {
		case err == nil:
		case ok && urlErr.Err == redirectError:
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
