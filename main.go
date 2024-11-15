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
	timeout := time.Duration(to) * time.Millisecond

	// Create a custom HTTP client with redirect handling
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}

	output := make(chan string)
	httpsURLs := make(chan string)

	// HTTPS workers
	var httpsWG sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		httpsWG.Add(1)

		go func() {
			defer httpsWG.Done()
			for url := range httpsURLs {
				output <- isListening(client, url)
			}
		}()
	}

	// Output worker
	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		defer outputWG.Done()
		for o := range output {
			fmt.Println(o)
		}
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

	// Close the input channel
	close(httpsURLs)

	// Check for errors reading stdin
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}
	// Wait until the output waitgroup is done
	outputWG.Wait()
}

func isListening(client *http.Client, rawURL string) string {
	// Add https if no protocol is present
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	// Try https first
	resp, err := client.Get(rawURL)
	if err != nil {
		// If https fails, try http
		httpURL := "http://" + rawURL[len("https://"):]
		resp, err = client.Get(httpURL)
		if err != nil {
			return fmt.Sprintf("Error: %s", err)
		}
	}

	defer resp.Body.Close() // Ensure the response body is closed
	urlFinal := resp.Request.URL.String()
	return urlFinal
}
