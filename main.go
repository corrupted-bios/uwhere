package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type probeArgs []string

func (p *probeArgs) Set(val string) error {
	*p = append(*p, val)
	return nil
}

func (p probeArgs) String() string {
	return strings.Join(p, ",")
}

func main() {

	// concurrency flag
	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "set the concurrency level")

	// timeout flag
	var to int
	flag.IntVar(&to, "t", 10000, "timeout (milliseconds)")

	flag.Parse()

	// make an actual time.Duration out of the timeout
	timeout := time.Duration(to * 1000000)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = timeout

	output := make(chan string)

	httpsURLs := make(chan string)

	// HTTPS workers
	var httpsWG sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		httpsWG.Add(1)

		go func() {
			for url := range httpsURLs {
				output <- isListening(url)
			}
			defer httpsWG.Done()
		}()
	}

	// Output worker
	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		for o := range output {
			fmt.Println(o)
		}
		outputWG.Done()
	}()

	// Close the output channel when the HTTP workers are done
	go func() {
		httpsWG.Wait()

		close(output)
	}()

	// accept domains on stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		domain := strings.ToLower(sc.Text())
		httpsURLs <- domain
	}

	// once we've sent all the URLs off we can close the
	// input/httpsURLs channel. The workers will finish what they're
	// doing and then call 'Done' on the WaitGroup
	close(httpsURLs)

	// check there were no errors reading stdin (unlikely)
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}
	// Wait until the output waitgroup is done
	outputWG.Wait()
}

// func isListening(client *http.Client, url, method string) bool {
func isListening(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		return ""
	}
	resp.Close = true
	urlFinal := resp.Request.URL.String()
	return urlFinal
}
