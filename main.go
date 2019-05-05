package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"strings"
)

const(
	version string = "0.1"
)

func (w *Website) String() string {
	var sb strings.Builder
	sb.WriteString(w.Url)
	sb.WriteString(" links the following websites:\n")
	for _, link := range w.LinkedSites {
		sb.WriteString("  ")
		sb.WriteString(link.Url)
	}
	return sb.String()
}

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
	if err != nil {
		panic(err)
	}

	versionArg, err := arguments.Bool("--version")
	if err == nil && versionArg {
		fmt.Printf("gora version %s", version)
		return
	}

	url, err := arguments.String("<url>")
	if err != nil {
		fmt.Println(usage)
		return
	}


	fmt.Printf("Starting from: %v\n", url)
	limit, _ := arguments.Int("--max")
	timeout, _ := arguments.Int("--timeout")

	if limit == 0 && timeout == 0 {
		fmt.Println("At least one of the options --max and --timeout needs to be set to a positive integer.")
		return
	}

	web := BuildGraph(url, limit, timeout, func(website *Website) {
		fmt.Printf("Processed: %s\n", website.Url)
	})
	fmt.Printf("The Web, starting from %s\n\n", url)
	for _, website := range web {
		fmt.Println(website)
	}

}