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
	"sync"
	"text/tabwriter"
	"time"

	"github.com/carlmjohnson/get-headers/prettyprint"
)

// Flag variables (set in flags.go on init)
var (
	etag       string
	gzip       bool
	ignoreBody bool
)

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
	client := http.Client{
		CheckRedirect: checkRedirect,
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	for _, url := range flag.Args() {

		req, err := http.NewRequest("GET", url, nil)
		die(err)

		req.Header = map[string][]string{}

		if gzip {
			req.Header["Accept-Encoding"] = []string{"gzip, deflate"}
		}

		if etag != "" {
			req.Header["If-None-Match"] = []string{etag}
		}

		start := time.Now()
		resp, err := client.Do(req)
		duration := time.Since(start)

		var (
			n  int64
			wg sync.WaitGroup
		)
		// Ignore the error if it's just our errRedirect
		switch urlErr, ok := err.(*netURL.Error); {
		case err == nil:
			if ignoreBody {
				break
			}

			wg.Add(1)
			go func() {
				// Copying to /dev/null just to make sure this is real
				n, err = io.Copy(ioutil.Discard, resp.Body)
				duration = time.Since(start)
				die(err)
				wg.Done()
			}()
		case ok && urlErr.Err == errRedirect:
		default:
			die(err)
		}

		fmt.Println("GET", url)
		fmt.Println(resp.Proto, resp.Status)
		fmt.Println()
		fmt.Println(prettyprint.ResponseHeader(resp.Header))
		wg.Wait()
		die(resp.Body.Close())

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "Time\t%s\n", prettyprint.Duration(duration))
		if n != 0 {
			fmt.Fprintf(tw, "Content length\t%s\n", prettyprint.Size(n))
			bps := prettyprint.Size(float64(n) / duration.Seconds())
			fmt.Fprintf(tw, "Speed\t%s/s\n", bps)
		}
		tw.Flush()
	}
}
