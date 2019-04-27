package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `gora, the web explorer. Starting from a given URL, crawls the web and creates a graph model of connected sites.

Usage:
  gora [options] <url>
  gora -h | --help | --version

Options:
  -h --help        Show this screen.
  --version        Show version information.
  --timeout=<sec>  Stop crawling after <sec> seconds.
  --max=<numurls>  Stop crawling after fetching <numurls> URLs.`

	arguments, err := docopt.ParseDoc(usage)
	fmt.Println(arguments)
	fmt.Println(err)
}