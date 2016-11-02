// Package cli parses command line flags for get-headers and calls run.Main
package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/carlmjohnson/get-headers/run"
)

const usage = `Usage of get-headers:

get-headers [opts] <url>...
        Gets the URLs and prints their headers alphabetically.
        Repeated headers are printed with an asterisk.

`

// Run parses command line options, runs run.Main, and returns an os.Exit code.
func Run() int {
	gzip := flag.Bool("gzip", false, "Enable GZIP compression")
	flag.BoolVar(gzip, "g", false, "Shortcut for -gzip")
	ignoreBody := flag.Bool("ignore-body", false, "Ignore body of request; close connection after gettings the headers")
	flag.BoolVar(ignoreBody, "i", false, "Shortcut for -ignore-body")
	etag := flag.String("etag", "", "Set 'If-None-Match' header to etag value")
	cookie := flag.String("cookie",
		os.Getenv("GET_HEADERS_COOKIE"),
		"Set cookie header (overrides GET_HEADERS_COOKIE environmental variable)")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}

	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		return 2
	}

	// Normalize etag...
	if !strings.HasPrefix(*etag, `"`) {
		*etag = fmt.Sprintf(`"%s"`, *etag)
	}

	if err := run.Main(*cookie, *etag, *gzip, *ignoreBody, flag.Args()...); err != nil {
		fmt.Fprintf(os.Stderr, "get-headers: %v\n", err)
		return 1
	}
	return 0
}
