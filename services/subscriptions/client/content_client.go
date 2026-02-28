package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type ContentClient struct {
	client  *http.Client
	cb      *gobreaker.CircuitBreaker
	baseURL string
}

func NewContentClient(baseURL string) *ContentClient {
	// 2.7.1 Konfiguracija HTTP Klijenta + 2.7.2 Timeout na nivou zahteva
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   3 * time.Second,
		Transport: tr,
	}

	// 2.7.4 Circuit Breaker
	cbSettings := gobreaker.Settings{
		Name:        "SubscriptionsToContent",
		MaxRequests: 3,
		Interval:    5 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker '%s' changed state from %s to %s", name, from, to)
		},
	}

	return &ContentClient{
		client:  client,
		cb:      gobreaker.NewCircuitBreaker(cbSettings),
		baseURL: baseURL,
	}
}

// CheckArtistExists proverava sinhrono da li umetnik postoji (Zahtev 2.5)
func (c *ContentClient) CheckArtistExists(artistID string) (bool, error) {
	url := fmt.Sprintf("%s/artists/%s", c.baseURL, artistID)
	log.Printf("[DEBUG] Subscriptions calling Content: %s", url)

	result, err := c.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return false, err
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return true, nil
		} else if resp.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			// 2.7.3 Fallback logika
			return c.fallbackResponse(artistID, err)
		}
		return false, err
	}

	return result.(bool), nil
}

func (c *ContentClient) fallbackResponse(targetID string, err error) (bool, error) {
	log.Printf("[FALLBACK] Could not verify target %s due to: %v. Fail-open for UX.", targetID, err)
	return true, nil
}
