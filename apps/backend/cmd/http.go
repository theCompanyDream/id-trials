package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/theCompanyDream/id-trials/apps/backend/models"
)

func TestApi(config *models.CmdConfig) {
	var wg sync.WaitGroup

	endpoints := []struct {
		name     string
		endpoint string
	}{
		{"ULID", "/api/ulid"},
		{"UUID", "/api/uuid4"},
		{"KSUID", "/api/ksuid"},
		{"CUID", "/api/cuid"},
		{"NanoID", "/api/nano"},
		{"Snowflake", "/api/snow"},
	}

	fmt.Printf("Load testing %d requests per endpoint across %d endpoints...\n",
		config.RecordsPerTable, len(endpoints))
	fmt.Printf("Concurrency: %d requests per endpoint\n", config.ConcurrentReqs)

	start := time.Now()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: config.RequestTimeout,
	}

	for _, ep := range endpoints {
		wg.Add(1)
		go func(name, endpoint string) {
			defer wg.Done()

			stats := loadTestEndpoint(client, config, name, endpoint)

			fmt.Printf("✅ %s: %d requests in %v (avg: %.2fms, success: %.1f%%)\n",
				name,
				stats.TotalRequests,
				stats.Duration,
				stats.AvgDuration,
				stats.SuccessRate,
			)
		}(ep.name, ep.endpoint)
	}

	wg.Wait()
	totalDuration := time.Since(start)

	fmt.Printf("\n🎉 Total: %d requests across all endpoints in %v\n",
		config.RecordsPerTable*len(endpoints), totalDuration)
}

type EndpointStats struct {
	TotalRequests int
	SuccessCount  int64
	ErrorCount    int64
	Duration      time.Duration
	AvgDuration   float64
	SuccessRate   float64
}

func loadTestEndpoint(client *http.Client, config *models.CmdConfig, name, endpoint string) EndpointStats {
	var (
		successCount  int64
		errorCount    int64
		totalDuration int64 // microseconds
		wg            sync.WaitGroup
		semaphore     = make(chan struct{}, config.ConcurrentReqs)
	)

	start := time.Now()

	for i := 0; i < config.RecordsPerTable; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire

		go func(index int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release

			userName := gofakeit.Username()
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			email := fmt.Sprintf("%c%s@%s.com", firstName[0], lastName, gofakeit.Company())
			department := &gofakeit.Job().Title

			// Generate fake user data
			user := models.UserInput{
				UserName:   &userName,
				FirstName:  &firstName,
				LastName:   &lastName,
				Email:      &email,
				Department: department,
			}

			// Make HTTP request
			reqStart := time.Now()
			err := createUser(client, config.BaseURL+endpoint, user)
			reqDuration := time.Since(reqStart)

			// Track stats
			atomic.AddInt64(&totalDuration, reqDuration.Microseconds())

			if err != nil {
				atomic.AddInt64(&errorCount, 1)
				if index%100 == 0 { // Don't spam errors
					fmt.Printf("  ⚠️  %s error: %v\n", name, err)
				}
			} else {
				atomic.AddInt64(&successCount, 1)
			}

			// Progress update
			if (index+1)%config.BatchSize == 0 {
				fmt.Printf("  %s progress: %d/%d requests\n", name, index+1, config.RecordsPerTable)
			}

			// Rate limiting
			if config.DelayBetweenReqs > 0 {
				time.Sleep(config.DelayBetweenReqs)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	// Calculate stats
	totalReqs := int(successCount + errorCount)
	avgDuration := float64(totalDuration) / float64(totalReqs) / 1000.0 // Convert to ms
	successRate := float64(successCount) / float64(totalReqs) * 100.0

	return EndpointStats{
		TotalRequests: totalReqs,
		SuccessCount:  successCount,
		ErrorCount:    errorCount,
		Duration:      duration,
		AvgDuration:   avgDuration,
		SuccessRate:   successRate,
	}
}

func createUser(client *http.Client, url string, user models.UserInput) error {
	// Marshal user to JSON
	payload, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("request creation error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Drain body to reuse connection
	io.Copy(io.Discard, resp.Body)

	return nil
}
