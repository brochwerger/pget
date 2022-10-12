package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	var nthreads = flag.Int("n", 1, "Number of parallel threads")
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: pget -n <number of getter threads> url")
		os.Exit(-1)
	}

	url := flag.Args()[0]
	ch := make(chan string)

	for thread := 0; thread < *nthreads; thread++ {
		go fetch(thread, url, ch) // start a goroutine
	}
	for thread := 0; thread < *nthreads; thread++ {
		fmt.Println(<-ch) // receive from channel ch
	}

	fmt.Println("Done")
}

func fetch(id int, url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}
	defer resp.Body.Close() // don't leak resources

	if resp.StatusCode != 200 {
		ch <- fmt.Sprintf("Error downloading %s: %v", url, err)
	}

	// Read the response body as a byte slice
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	mimeType := http.DetectContentType(bytes)

	// Determined the response type and print message accordingly
	var suffix string
	if strings.Contains(mimeType, "text") {
		suffix = string(bytes)
	} else {
		nbytes := len(bytes)
		suffix = fmt.Sprintf(" %dKB downloaded (%s)", nbytes/1024, mimeType)
	}
	secs := time.Since(start)
	ch <- fmt.Sprintf("Thread %v finished in %v: %s", id, secs, suffix)
}
