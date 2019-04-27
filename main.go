package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"net/http"
)

const(
	version string = "0.1"
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
	if err != nil {
		panic(err)
	}

	versionArg, err := arguments.Bool("--version")
	if err == nil && versionArg {
		fmt.Printf("gora version %s", version)
	}

	url, err := arguments.String("<url>")
	if err != nil {
		fmt.Println(usage)
		return
	}

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Could not get %s, error: %s", url, err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Could not read body of %s, error: %s", url, err)
	}

	fmt.Print(string(body))
}