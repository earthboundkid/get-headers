// Package cli parses command line flags for get-headers and calls run.Main
package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/carlmjohnson/exitcode"
	"github.com/carlmjohnson/flagext"
	"github.com/carlmjohnson/get-headers/run"
	"github.com/carlmjohnson/versioninfo"
)

const usage = `Usage of get-headers %s:

get-headers [opts] <url>...
        Gets the URLs and prints their headers alphabetically.
        Repeated headers are printed with an asterisk.

Options may be set as GET_HEADERS prefixed environment variables.

`

// Run parses command line options, runs run.Main, and returns an os.Exit code.
func Run() int {
	gzip := flag.Bool("gzip", false, "Enable GZIP compression")
	flag.BoolVar(gzip, "g", false, "Shortcut for -gzip")
	ignoreBody := flag.Bool("ignore-body", false, "Ignore body of request; close connection after gettings the headers")
	flag.BoolVar(ignoreBody, "i", false, "Shortcut for -ignore-body")
	etag := ""
	flag.Func("etag", "Set 'If-None-Match' header to etag value", func(s string) error {
		// Normalize etag...
		if !strings.HasPrefix(etag, `"`) {
			etag = fmt.Sprintf(`"%s"`, etag)
		}
		return nil
	})
	cookie := flag.String("cookie", "", "Set cookie header")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, versioninfo.Version)
		flag.PrintDefaults()
	}

	flag.Parse()
	flagext.ParseEnv(nil, "get-headers")
	flagext.MustHaveArgs(nil, 1, -1)

	return exitcode.Get(run.Main(*cookie, etag, *gzip, *ignoreBody, flag.Args()...))
}
