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

func init() {
	flag.BoolVar(&gzip, "gzip", false, "Enable GZIP compression")
	flag.BoolVar(&gzip, "g", false, "Shortcut for -gzip")
	flag.BoolVar(&ignoreBody, "ignore-body", false, "Ignore body of request; close connection after gettings the headers")
	flag.BoolVar(&ignoreBody, "i", false, "Shortcut for -ignore-body")
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
