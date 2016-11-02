// Package run is the core logic of get-headers
package run

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	netURL "net/url"
	"os"
	"text/tabwriter"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/carlmjohnson/get-headers/prettyprint"
)

// Sentinal error to let us know if we're ignoring a redirect
var errRedirect = errors.New("redirected")

// Don't follow redirects
func checkRedirect(req *http.Request, via []*http.Request) error {
	return errRedirect
}

// client for all http requests
var client = http.Client{
	CheckRedirect: checkRedirect,
	Transport: &http.Transport{
		DisableCompression: true,
	},
}

// Main takes a list of urls and request parameters, then fetches the URLs and
// outputs the headers to stdout
func Main(cookie, etag string, gzip, ignoreBody bool, urls ...string) error {
	for nurl, url := range urls {
		// Separate subsequent lookups with newline
		if nurl > 0 {
			fmt.Println()
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		if gzip {
			req.Header.Add("Accept-Encoding", "gzip, deflate")
		}

		if etag != "" {
			req.Header.Add("If-None-Match", etag)
		}

		if cookie != "" {
			req.Header.Add("Cookie", cookie)
		}

		start := time.Now()
		resp, err := client.Do(req)
		duration := time.Since(start)

		var (
			n  int64
			eg errgroup.Group
		)
		// Ignore the error if it's just our errRedirect
		switch urlErr, ok := err.(*netURL.Error); {
		case err == nil:
			if ignoreBody {
				break
			}

			eg.Go(func() error {
				// Copying to /dev/null just to make sure this is real
				n, err = io.Copy(ioutil.Discard, resp.Body)
				duration = time.Since(start)
				if err != nil {
					return err
				}
				return nil
			})
		case ok && urlErr.Err == errRedirect:
		default:
			return err
		}

		fmt.Println("GET", url)
		fmt.Println(resp.Proto, resp.Status)
		fmt.Println()
		fmt.Println(prettyprint.ResponseHeader(resp.Header))

		if err := eg.Wait(); err != nil {
			return err
		}

		if err := resp.Body.Close(); err != nil {
			return err
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "Time\t%s\n", prettyprint.Duration(duration))
		if n != 0 {
			fmt.Fprintf(tw, "Content length\t%s\n", prettyprint.Size(n))
			bps := prettyprint.Size(float64(n) / duration.Seconds())
			fmt.Fprintf(tw, "Speed\t%s/s\n", bps)
		}
		if err := tw.Flush(); err != nil {
			return err
		}
	}

	return nil
}
