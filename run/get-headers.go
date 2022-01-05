// Package run is the core logic of get-headers
package run

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"text/tabwriter"
	"time"

	"github.com/carlmjohnson/get-headers/prettyprint"
	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/versioninfo"
)

func IPDialer() (*net.Addr, *http.Client) {
	var ip net.Addr
	t := &http.Transport{
		DisableCompression: true,
		ForceAttemptHTTP2:  true,
	}
	t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := net.Dial(network, addr)
		if conn != nil {
			ip = conn.RemoteAddr()
		}
		return conn, err
	}
	cl := &http.Client{
		CheckRedirect: requests.NoFollow,
		Transport:     t,
	}
	return &ip, cl
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

	builder := requests.URL(url)
	builder.UserAgent(getUserAgent())
	if gzip {
		builder.Header("Accept-Encoding", "gzip, deflate")
	}

	if etag != "" {
		builder.Header("If-None-Match", etag)
	}

	if cookie != "" {
		builder.Header("Cookie", cookie)
	}

	ip, newClient := IPDialer()
	builder.Client(newClient)

	var (
		size             int64
		start            time.Time
		duration         time.Duration
		err              error
		printheadersDone = make(chan struct{})
	)
	builder.AddValidator(func(res *http.Response) error {
		go func() {
			fmt.Println("GET", url)
			if *ip != nil {
				fmt.Println("Via", *ip)
			}
			fmt.Println(res.Proto, res.Status)
			fmt.Println()
			fmt.Println(prettyprint.ResponseHeader(res.Header))
			close(printheadersDone)
		}()
		return nil
	})
	if ignoreBody {
		builder.Handle(func(res *http.Response) error {
			duration = time.Since(start)
			return nil
		})
	} else {
		builder.Handle(func(res *http.Response) error {
			size, err = io.Copy(io.Discard, res.Body)
			duration = time.Since(start)
			return err
		})
	}

	start = time.Now()
	if err = builder.Fetch(ctx); err != nil {
		return err
	}
	<-printheadersDone
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Time\t%s\n", prettyprint.Duration(duration))
	if size != 0 {
		fmt.Fprintf(tw, "Content length\t%s\n", prettyprint.Size(size))
		bps := prettyprint.Size(float64(size) / duration.Seconds())
		fmt.Fprintf(tw, "Speed\t%s/s\n", bps)
	}
	if err := tw.Flush(); err != nil {
		return err
	}

	return nil
}

var userAgent string

func getUserAgent() string {
	if userAgent != "" {
		return userAgent
	}
	userAgent = fmt.Sprintf("get-headers/%s", versioninfo.Version)
	return userAgent
}
