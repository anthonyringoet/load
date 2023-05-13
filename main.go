package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

var (
	times     []time.Duration
	successes int
	failures  int
	running   bool
	wg        sync.WaitGroup
)

func sendRequests(url string, verbose bool) {
	defer wg.Done()

	// Create an HTTP client
	client := &http.Client{}

	for running {
		start := time.Now()

		// Create a new request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Set the User-Agent header
		req.Header.Set("User-Agent", "github.com/anthonyringoet/load/0.0.1")

		// Send the request using the same client
		resp, err := client.Do(req)
		if err != nil {
			failures++
			if verbose {
				fmt.Printf("Failed to send request: %s\n", err)
			}
			continue
		}

		elapsed := time.Since(start)
		times = append(times, elapsed)

		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			successes++
			if verbose {
				fmt.Printf("✅ success. statuscode=%d http_version=%s duration=%dms\n", resp.StatusCode, resp.Proto, elapsed.Milliseconds())
			}
		} else {
			failures++
			if verbose {
				fmt.Printf("❌ fail. statuscode=%d http_version=%s duration=%dms\n", resp.StatusCode, resp.Proto, elapsed.Milliseconds())
			}
		}
	}
}

// calculateStats takes a slice of time.Duration values representing latencies,
// sorts them, and then calculates and prints various statistics:
//
//   - Min latency: The smallest latency in the slice.
//   - Max latency: The largest latency in the slice.
//   - Median latency: The middle value in the sorted slice. If the slice has an even
//     number of values, this is the lower of the two middle values.
//   - 90th percentile latency: The value below which 90% of the latencies fall.
//   - 95th percentile latency: The value below which 95% of the latencies fall.
//   - 99th percentile latency: The value below which 99% of the latencies fall.
//   - Average latency: The sum of all latencies divided by the number of latencies.
//
// The function does not return a value; instead, it prints the calculated statistics
// to standard output.
//
// The input slice is sorted in-place, so the order of values in the slice will be
// changed by this function.
func calculateStats(times []time.Duration) {
	if len(times) == 0 {
		return
	}

	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})

	min := times[0]
	max := times[len(times)-1]

	median := times[len(times)/2]

	p90 := times[int(float64(len(times))*0.9)]
	p95 := times[int(float64(len(times))*0.95)]
	p99 := times[int(float64(len(times))*0.99)]

	// calculate average
	var total time.Duration
	for _, t := range times {
		total += t
	}

	average := total / time.Duration(len(times))

	fmt.Printf("Min latency: %s\n", min)
	fmt.Printf("Max latency: %s\n", max)
	fmt.Printf("Median latency: %s\n", median)
	fmt.Printf("90th percentile latency: %s\n", p90)
	fmt.Printf("95th percentile latency: %s\n", p95)
	fmt.Printf("99th percentile latency: %s\n", p99)
	fmt.Printf("Average latency: %s\n", average)
}

func logProgress(duration *time.Duration) {
	// simple progress updates
	for i := 1; i <= int((*duration).Seconds()); i++ {
		time.Sleep(1 * time.Second)
		if running {
			fmt.Printf("%d/%d seconds elapsed\n", i, int((*duration).Seconds()))
		}
	}
}

func main() {
	url := flag.String("url", "http://example.com", "URL to load test")
	workers := flag.Int("workers", 1, "Number of workers")
	duration := flag.Duration("duration", 10*time.Second, "Test duration")
	verbose := flag.Bool("verbose", false, "Verbose output")
	flag.Parse()

	if *duration < 0 {
		fmt.Println("Duration cannot be negative. Defaulting to 10 seconds")
		*duration = 10 * time.Second
	}
	if *workers < 0 {
		fmt.Println("Workers cannot be negative. Defaulting to 1 worker")
		*workers = 1
	}

	fmt.Printf("Load testing %s with %d workers for %s\n\n", *url, *workers, *duration)

	running = true

	if !*verbose {
		go logProgress(duration)
	}

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go sendRequests(*url, *verbose)
	}

	time.Sleep(*duration)

	running = false
	wg.Wait()

	fmt.Printf("\n2xx requests: %d, non 2xx requests: %d\n", successes, failures)
	fmt.Printf("Total requests: %d\n", successes+failures)
	fmt.Printf("Avg requests/worker: %d\n\n", (successes+failures)/(*workers))
	calculateStats(times)
	fmt.Println("\nLoad testing finished")
}
