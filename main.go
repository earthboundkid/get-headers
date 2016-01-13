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
	"time"

	"github.com/carlmjohnson/get-headers/prettyprint"
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
	transport := &http.Transport{
		DisableCompression: *gzip,
	}
	client := http.Client{
		CheckRedirect: checkRedirect,
		Transport:     transport,
	}
	for _, url := range flag.Args() {

		req, err := http.NewRequest("GET", url, nil)
		die(err)
		if *gzip {
			req.Header = map[string][]string{
				"Accept-Encoding": {"gzip, deflate"},
			}
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
		fmt.Println(resp.Proto, resp.Status, "\n")
		fmt.Println(prettyprint.ResponseHeader(resp.Header))
		wg.Wait()
		die(resp.Body.Close())
		fmt.Println("Time           ", humanizeDuration(duration))
		if n != 0 {
			fmt.Println("Content length ", humanizeByteSize(n))
			bps := int64(float64(n) / duration.Seconds())
			fmt.Printf("Speed           %s/s\n", humanizeByteSize(bps))
		}
		fmt.Println()
	}
}
