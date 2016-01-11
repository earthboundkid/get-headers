// get-headers prints the headers from GET requesting a URL
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	netURL "net/url"
	"os"
	"time"

	"github.com/carlmjohnson/get-headers/prettyprint"
)

const usage = `Usage of get-headers:

get-headers [opts] <url>...
       	Gets the URLs and prints their headers alphabetically.
       	Repeated headers are printed with an asterisk.

`

var (
	gzip = flag.Bool("gzip", false, "Enable GZIP compression")
)

func init() {
	flag.BoolVar(gzip, "g", false, "Shortcut for -gzip")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
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

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "get-time: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	transport := &http.Transport{
		DisableCompression: *gzip,
	}
	client := http.Client{
		CheckRedirect: checkRedirect,
		Transport:     transport,
	}
	for _, url := range flag.Args() {
		var n int64
		start := time.Now()

		req, err := http.NewRequest("GET", url, nil)
		die(err)
		if *gzip {
			req.Header = map[string][]string{
				"Accept-Encoding": {"gzip, deflate"},
			}
		}
		resp, err := client.Do(req)

		// Ignore the error if it's just our errRedirect
		switch urlErr, ok := err.(*netURL.Error); {
		case err == nil:
			// Copying to /dev/null just to make sure this is real
			n, err = io.Copy(ioutil.Discard, resp.Body)
			die(err)
		case ok && urlErr.Err == errRedirect:
		default:
			die(err)
		}
		duration := time.Since(start)
		die(resp.Body.Close())

		fmt.Println("GET", url)
		fmt.Println(resp.Proto, resp.Status, "\n")
		fmt.Println("Time           ", humanizeDuration(duration))
		if n != 0 {
			fmt.Println("Content length ", humanizeByteSize(n))
			bps := int64(float64(n) / duration.Seconds())
			fmt.Printf("Speed           %s/s\n", humanizeByteSize(bps))
		}
		fmt.Println()
		fmt.Println(prettyprint.ResponseHeader(resp.Header))

	}
}
