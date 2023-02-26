package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	timeout     int
	concurrency int
	verbose     bool
	method      string
)

// Visit It holds the URL and the size of the response size
type Visit struct {
	Url      string
	BodySize int
	Error    error
}

func main() {
	// command line flags
	flag.IntVar(&timeout, "t", 30, "request timeout for the http request in seconds")
	flag.IntVar(&concurrency, "c", 1, "number of maximum of concurrent requests")
	flag.BoolVar(&verbose, "v", false, "if provided, enables the verbose mode, it outputs request errors")
	flag.StringVar(&method, "m", "GET", "http method to make the HTTP call")

	var visits []*Visit

	flag.Parse()
	args := flag.Args()

	// validate integer params
	if timeout < 1 || timeout > 600 {
		fmt.Println("-t timeout param must be in the range [1, 600]")
		os.Exit(1)
	}

	// validate integer params
	if concurrency < 1 || concurrency > 100 {
		fmt.Println("-c concurrency param must be in the range [1, 100]")
		os.Exit(2)
	}

	// validate http methods, not really necessary, but I'd like to verify before create the request
	if !contains([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}, strings.ToUpper(method)) {
		fmt.Println("-m http unsupported http method, use valid values [GET, HEAD, POST, PUT, PATCH, DELETE, " +
			"CONNECT, OPTIONS, TRACE]")
		os.Exit(3)
	}

	if len(args) > 0 { // if at least one parameter is provided

		if concurrency > 1 { // concurrent processing
			visitQ := make(chan Visit, concurrency)

			for _, urlAddr := range args {
				go visitUrl(nil, urlAddr, method, timeout, true, verbose, visitQ)
			}

			for i := 0; i < len(args); i++ {
				visit := <-visitQ
				visits = append(visits, &visit)
			}

		} else { // serial processing, useful when you have network issues
			for _, urlAddr := range args {
				visit := visitUrl(nil, urlAddr, method, timeout, false, verbose, nil) // no need to use the WaitGroup here
				if visit.Error != nil {
					if verbose {
						fmt.Println(visit.Error.Error())
					}
				} else {
					visits = append(visits, visit)
				}
			}
		}

		// sort them out
		sortVisits(visits, true)

		// print response to terminal
		for _, visit := range visits {
			if visit.Error != nil {
				fmt.Printf("%s %d bytes (%s)\n", visit.Url, visit.BodySize, visit.Error)
			} else {
				fmt.Printf("%s %d bytes\n", visit.Url, visit.BodySize)
			}
		}

	} else {
		fmt.Println("For help run: \n$ wsc -h")
		fmt.Println("")
		fmt.Println("Example with default values:")
		fmt.Println("$ wsc https://...de/ https://...com/")
		fmt.Println("https://...de/ 14953 bytes\nhttps://...com/ 359600 bytes")
		fmt.Println("")
		fmt.Println("Example with request timeout in seconds:")
		fmt.Println("$ wsc -t 5 https://...de/ https://...com/")
		fmt.Println("https://...de/ 14953 bytes\nhttps://...com/ 359600 bytes")
		fmt.Println("")
		fmt.Println("Example with max concurrent requests:")
		fmt.Println("$ wsc -c 2 https://...de/ https://...com/")
		fmt.Println("https://...de/ 14953 bytes\nhttps://...com/ 359600 bytes")
		fmt.Println("")
		fmt.Println("Example with POST method request:")
		fmt.Println("$ wsc -c 2 -m POST https://...de/ https://...com/")
		fmt.Println("https://...de/ 1453 bytes\nhttps://...com/ 5960 bytes")
		fmt.Println("")
		fmt.Println("Example with verbose mode enabled:")
		fmt.Println("$ wsc -v -c 2 -t 1 -m POST https://...de/ https://...com/ https://....com.br")
		fmt.Println("https://....com.br 0 bytes (Post \"https://....com.br\": context deadline exceeded (Client.Timeout exceeded while awaiting headers))\nhttps://....de/ 936 bytes\nhttps://...com/ 1861 bytes")
		fmt.Println("")
	}
}

// visitUrl Fetches the URL data
func visitUrl(client *http.Client, urlAddress, method string, timeout int, concurrent, verbose bool, visitQ chan Visit) *Visit {
	visit := &Visit{Url: urlAddress}

	// in case of concurrent processing, feedback
	defer func(visit *Visit, concurrent bool, visitQ chan Visit) {
		// in case the concurrent execution
		if concurrent {
			visitQ <- *visit
		}
	}(
		visit,
		concurrent,
		visitQ,
	)

	// let fetch the content
	if client == nil {
		client = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}

	req, err := http.NewRequest(method, urlAddress, nil)
	if err != nil {
		if verbose {
			fmt.Printf("ERROR: cannot create a valid request - ignoring invalid URL '%s', setting body size to zero\n", urlAddress)
		}
		visit.Error = err

		return visit
	}

	resp, err := client.Do(req)
	if err != nil {
		visit.Error = err

		return visit
	}

	// Usually resp.ContentLength is unreliable thus
	// we need to use `httputil.DumpResponse` to extract
	// all information we need
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		visit.Error = err
		return visit
	}

	visit.BodySize = len(b)

	return visit
}

// contains Used for verifying if a http method is in the valid range.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// sortVisits Sort all visits by the Visit.BodySize
func sortVisits(visits []*Visit, asc bool) {
	// order slice taking advantage of the `sort.Interface`
	if asc {
		sort.SliceStable(visits, func(i, j int) bool {
			return visits[i].BodySize < visits[j].BodySize // ascending order
		})
	} else {
		sort.SliceStable(visits, func(i, j int) bool {
			return visits[i].BodySize > visits[j].BodySize // descending order
		})
	}
}
