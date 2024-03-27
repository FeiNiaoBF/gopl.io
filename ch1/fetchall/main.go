// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 17.
//!+

// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start1 := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		// if !strings.HasPrefix(url, "http://") {
		// 	url = "http://" + url
		// }
		go fetch(url, ch) // start a goroutine
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start1).Seconds())

	start2 := time.Now()
	for _, url := range os.Args[1:] {
		fetchBuffer(url)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start2).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("Test1: %.2fs  %7d  %s", secs, nbytes, url)
}

//!-

func fetchBuffer(url string) {
	start := time.Now()
	// if !strings.HasPrefix(url, "http://") {
	// 	url = "http://" + url
	// }
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	result := fmt.Sprintf("Test2: %.2fs  %7d  %s", secs, nbytes, url)
	fmt.Println(result)
}
