package main

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "test/prot/signuppb" // Update this path based on your actual module name
)

type result struct {
	duration time.Duration
	success  bool
}

func main() {
	address := "localhost:8080" // gRPC server address
	initialRequests := 100
	increment := 100
	steps := 50
	concurrency := 100

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("Failed to connect: %v", err))
	}
	defer conn.Close()
	client := pb.NewAuthServiceClient(conn)

	for step := 1; step <= steps; step++ {
		totalRequests := initialRequests + (step-1)*increment
		fmt.Printf("\nStep %d: Sending %d requests with concurrency %d...\n", step, totalRequests, concurrency)

		startTime := time.Now()
		results := runGRPCLoadTest(client, totalRequests, concurrency)
		totalDuration := time.Since(startTime)

		var durations []time.Duration
		successCount := 0
		for _, res := range results {
			if res.success {
				successCount++
			}
			durations = append(durations, res.duration)
		}
		failCount := len(results) - successCount

		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
		p50 := durations[len(durations)*50/100]
		p95 := durations[len(durations)*95/100]
		p99 := durations[len(durations)*99/100]
		min := durations[0]
		max := durations[len(durations)-1]
		rps := float64(totalRequests) / totalDuration.Seconds()

		fmt.Printf("Latency stats (ms):\n")
		fmt.Printf("P50: %.2f | P95: %.2f | P99: %.2f | Min: %.2f | Max: %.2f\n",
			toMs(p50), toMs(p95), toMs(p99), toMs(min), toMs(max))
		fmt.Printf("Total duration: %v | Requests/sec: %.2f\n", totalDuration, rps)
		fmt.Printf("Success: %d | Failed: %d\n", successCount, failCount)

		time.Sleep(1 * time.Second)
	}

	fmt.Println("\nBenchmark test completed.")
}

func runGRPCLoadTest(client pb.AuthServiceClient, total, concurrency int) []result {
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	results := make([]result, total)

	for i := 0; i < total; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			start := time.Now()
			_, err := client.Signup(ctx, &pb.SignupRequest{
				Username: fmt.Sprintf("user_%d", time.Now().UnixNano()),
				Password: "123",
			})
			elapsed := time.Since(start)

			results[i] = result{
				duration: elapsed,
				success:  err == nil,
			}
		}(i)
	}

	wg.Wait()
	return results
}

func toMs(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000.0
}
