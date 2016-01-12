package main

import (
	"flag"
	"fmt"
	"os"
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
