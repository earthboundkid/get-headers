// Package run is the core logic of get-headers
package run

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"text/tabwriter"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/carlmjohnson/get-headers/prettyprint"
)

// Don't follow redirects
func checkRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
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
	for i, url := range urls {
		// Separate subsequent lookups with newline
		if i > 0 {
			fmt.Println()
		}
		if err := getHeaders(cookie, etag, gzip, ignoreBody, url); err != nil {
			return err
		}
	}
	return nil
}

func getHeaders(cookie, etag string, gzip, ignoreBody bool, url string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

	if err != nil {
		return err
	}

	var n int64

	if !ignoreBody {
		eg.Go(func() error {
			// Copying to /dev/null just to make sure this is real
			n, err = io.Copy(io.Discard, resp.Body)
			duration = time.Since(start)
			if err != nil {
				return err
			}
			return nil
		})
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

	return nil
}
