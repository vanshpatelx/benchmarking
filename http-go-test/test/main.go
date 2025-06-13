// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"net/url"
// 	"sync"
// 	"time"
// 	"sort"
// )

// func main() {
// 	const concurrency = 50
// 	const totalRequests = 1000

// 	var wg sync.WaitGroup
// 	sem := make(chan struct{}, concurrency)
// 	timings := make([]time.Duration, totalRequests)

// 	for i := 0; i < totalRequests; i++ {
// 		wg.Add(1)
// 		sem <- struct{}{}

// 		go func(i int) {
// 			defer wg.Done()
// 			defer func() { <-sem }()

// 			data := url.Values{}
// 			data.Set("username", fmt.Sprintf("user_%d", i))
// 			data.Set("password", "123")

// 			start := time.Now()
// 			_, err := http.PostForm("http://localhost:8080/signup", data)
// 			timings[i] = time.Since(start)

// 			if err != nil {
// 				fmt.Printf("Request %d failed: %v\n", i, err)
// 			}
// 		}(i)
// 	}

// 	wg.Wait()
// 	fmt.Println("All done.")

// 	// Calculate latency stats
// 	sort.Slice(timings, func(i, j int) bool { return timings[i] < timings[j] })

// 	fmt.Printf("\nLatency stats (ms):\n")
// 	fmt.Printf("P50:  %v\n", timings[totalRequests*50/100])
// 	fmt.Printf("P95:  %v\n", timings[totalRequests*95/100])
// 	fmt.Printf("P99:  %v\n", timings[totalRequests*99/100])
// 	fmt.Printf("Max:  %v\n", timings[totalRequests-1])
// 	fmt.Printf("Min:  %v\n", timings[0])
// }

// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"net/url"
// 	"sort"
// 	"sync"
// 	"time"
// )

// func main() {
// 	baseURL := "http://localhost:8080/signup"
// 	initialRequests := 100
// 	increment := 100
// 	steps := 10
// 	concurrency := 50

// 	for step := 1; step <= steps; step++ {
// 		totalRequests := initialRequests + (step-1)*increment
// 		fmt.Printf("\nStep %d: Sending %d requests with concurrency %d...\n", step, totalRequests, concurrency)

// 		startTime := time.Now()
// 		timings := runLoadTest(baseURL, totalRequests, concurrency)
// 		totalDuration := time.Since(startTime)

// 		// Sort and compute stats
// 		sort.Slice(timings, func(i, j int) bool { return timings[i] < timings[j] })
// 		p50 := timings[len(timings)*50/100]
// 		p95 := timings[len(timings)*95/100]
// 		p99 := timings[len(timings)*99/100]
// 		min := timings[0]
// 		max := timings[len(timings)-1]
// 		rps := float64(totalRequests) / totalDuration.Seconds()

// 		fmt.Printf("Latency stats (ms):\n")
// 		fmt.Printf("P50: %v | P95: %v | P99: %v | Min: %v | Max: %v\n", p50, p95, p99, min, max)
// 		fmt.Printf("Total duration: %v | Requests/sec: %.2f\n", totalDuration, rps)

// 		time.Sleep(1 * time.Second)
// 	}

// 	fmt.Println("\nBenchmark test completed.")
// }

// func runLoadTest(endpoint string, total, concurrency int) []time.Duration {
// 	var wg sync.WaitGroup
// 	sem := make(chan struct{}, concurrency)
// 	timings := make([]time.Duration, total)

// 	for i := 0; i < total; i++ {
// 		wg.Add(1)
// 		sem <- struct{}{}

// 		go func(i int) {
// 			defer wg.Done()
// 			defer func() { <-sem }()

// 			form := url.Values{}
// 			form.Set("username", fmt.Sprintf("user_%d", time.Now().UnixNano()))
// 			form.Set("password", "123")

// 			start := time.Now()
// 			resp, err := http.PostForm(endpoint, form)
// 			timings[i] = time.Since(start)

// 			if err != nil {
// 				fmt.Printf("Request %d failed: %v\n", i, err)
// 				return
// 			}
// 			defer resp.Body.Close()
// 		}(i)
// 	}

// 	wg.Wait()
// 	return timings
// }


package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"
)

type result struct {
	duration time.Duration
	success  bool
}

func main() {
	baseURL := "http://localhost:8080/signup"
	initialRequests := 100
	increment := 100
	steps := 50
	concurrency := 100

	for step := 1; step <= steps; step++ {
		totalRequests := initialRequests + (step-1)*increment
		fmt.Printf("\nStep %d: Sending %d requests with concurrency %d...\n", step, totalRequests, concurrency)

		startTime := time.Now()
		results := runLoadTest(baseURL, totalRequests, concurrency)
		totalDuration := time.Since(startTime)

		// Filter durations and count success/failures
		var durations []time.Duration
		successCount := 0
		for _, res := range results {
			if res.success {
				successCount++
			}
			durations = append(durations, res.duration)
		}
		failCount := len(results) - successCount

		// Sort and compute latency percentiles
		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
		p50 := durations[len(durations)*50/100]
		p95 := durations[len(durations)*95/100]
		p99 := durations[len(durations)*99/100]
		min := durations[0]
		max := durations[len(durations)-1]
		rps := float64(totalRequests) / totalDuration.Seconds()

		fmt.Printf("Latency stats (ms):\n")
		fmt.Printf("P50: %v | P95: %v | P99: %v | Min: %v | Max: %v\n", toMs(p50), toMs(p95), toMs(p99), toMs(min), toMs(max))
		fmt.Printf("Total duration: %v | Requests/sec: %.2f\n", totalDuration, rps)
		fmt.Printf("Success: %d | Failed: %d\n", successCount, failCount)

		time.Sleep(1 * time.Second)
	}

	fmt.Println("\nBenchmark test completed.")
}

func runLoadTest(endpoint string, total, concurrency int) []result {
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	results := make([]result, total)

	for i := 0; i < total; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			form := url.Values{}
			form.Set("username", fmt.Sprintf("user_%d", time.Now().UnixNano()))
			form.Set("password", "123")

			start := time.Now()
			resp, err := http.PostForm(endpoint, form)
			elapsed := time.Since(start)

			success := err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300
			results[i] = result{duration: elapsed, success: success}

			if resp != nil {
				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()
	return results
}

// func toMs(d time.Duration) time.Duration {
// 	return time.Duration(float64(d) / float64(time.Millisecond))
// }


func toMs(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000.0
}
