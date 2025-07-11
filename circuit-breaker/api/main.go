package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/lmittmann/tint"
	"github.com/sony/gobreaker/v2"
)

var cb *gobreaker.CircuitBreaker[[]byte]
var circuitOpenedAt time.Time
var logger *slog.Logger

func init() {
	// Use tint for colorized terminal logs
	logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
	}))

	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		if counts.Requests == 0 {
			return false
		}
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}
	st.Timeout = 5 * time.Second
	st.OnStateChange = func(name string, from, to gobreaker.State) {
		logger.Warn("Circuit state changed",
			"name", name,
			"from", from.String(),
			"to", to.String(),
		)

		if to == gobreaker.StateOpen {
			circuitOpenedAt = time.Now()
		}

		if from == gobreaker.StateOpen && (to == gobreaker.StateHalfOpen || to == gobreaker.StateClosed) {
			duration := time.Since(circuitOpenedAt)
			logger.Info("Circuit was open for", "duration", duration)
		}
	}

	cb = gobreaker.NewCircuitBreaker[[]byte](st)
}

func Get(url string) ([]byte, error) {
	body, err := cb.Execute(func() ([]byte, error) {
		client := &http.Client{
			Timeout: 2 * time.Second,
		}

		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Check for logical error in JSON body
		var parsed map[string]any
		if err := json.Unmarshal(body, &parsed); err == nil {
			if errMsg, ok := parsed["error"].(string); ok && errMsg != "" {
				return nil, fmt.Errorf("logical error from DB: %s", errMsg)
			}
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
		}

		return body, nil
	})

	if err != nil {
		return nil, err
	}
	return body, nil
}

func main() {
	for i := 1; i <= 100; i++ {
		logger.Info("Sending request", "iteration", i)
		body, err := Get("http://localhost:8081/ping") // replace with real endpoint
		if err != nil {
			logger.Error("Request failed", "error", err)
		} else {
			logger.Info("Request success", "response", string(body))
		}
		time.Sleep(1 * time.Second)
	}
}
