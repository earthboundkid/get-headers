// get-headers prints the headers from GET requesting a URL
package main

import (
	"os"

	"github.com/carlmjohnson/get-headers/cli"
)

func main() {
	os.Exit(cli.Run())
}
