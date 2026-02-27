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
	// Use custom Transport to skip TLS verification for internal self-signed certs
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   3 * time.Second,
		Transport: tr,
	}

	// 2.7.4 Circuit Breaker
	cbSettings := gobreaker.Settings{
		Name:        "ContentService",
		MaxRequests: 3,                // broj requestova dozvoljenih kada je half-open
		Interval:    5 * time.Second,  // vreme čišćenja brojača
		Timeout:     10 * time.Second, // vreme čekanja pre prebacivanja u half-open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Tripuje (prelazi u Open) ako od minimum 3 zahteva, preko 60% padne
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

// CheckSongExists šalje GET zahtev prema Content servisu da proveri postojanje pesme
func (c *ContentClient) CheckSongExists(songID string) (bool, error) {
	url := fmt.Sprintf("%s/songs/%s", c.baseURL, songID)

	var lastErr error

	// 2.7.5 Retry mehanizam (max 3 pokušaja)
	for attempt := 1; attempt <= 3; attempt++ {
		// Circuit Breaker execute omotava poziv
		result, err := c.cb.Execute(func() (interface{}, error) {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return false, err
			}

			resp, err := c.client.Do(req)
			if err != nil {
				return false, err // greške usled timeout-a ili nedostupnosti mreže okidaju CB
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				return true, nil
			} else if resp.StatusCode == http.StatusNotFound {
				// pesma legit ne postoji (nije greška infrastrukture)
				return false, nil
			}

			// server je vratio grešku 5xx, okida CB
			return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		})

		if err == nil {
			return result.(bool), nil // Usašno ili potvrđeno da ne postoji
		}

		if errors.Is(err, gobreaker.ErrOpenState) {
			// 2.7.3 Fallback logiku okidamo odma ako je CB otvoren i nema svrhe retry-ovati
			return c.fallbackResponse(songID, attempt, err)
		}

		lastErr = err
		log.Printf("Attempt %d failed to check song (Content Service): %v", attempt, err)
		time.Sleep(1 * time.Second) // backoff
	}

	// 2.7.3 Fallback logika ako propadnu svi retry pokušaji
	return c.fallbackResponse(songID, 4, lastErr) // 'attempt' stavljen na 4 kako bi u logu izgledalo logicno nakon 3 pokušaja
}

func (c *ContentClient) fallbackResponse(songID string, attempt int, err error) (bool, error) {
	log.Printf("[FALLBACK] Could not verify song %s due to: %v. Suggesting default valid behaviour to avoid breaking rating logic", songID, err)
	// Fallback akcija - Ako ne možemo da pričamo sa Content Servisom, pretpostavićemo privremeno
	// da pesma postoji kako bi klijent mogao da izlista i upiše ocenu bez problema.
	// Druga opcija je vratiti error. Ovde biramo fail-open fallback.
	return true, nil
}
