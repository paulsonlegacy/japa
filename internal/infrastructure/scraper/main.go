package scraper

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)


type Scraper interface {
	Scrape() error
}

// MultiScraper has access to every scraper functions
type MultiScraper struct {
	Scrapers     []Scraper // List of scrapers in priority order
	Logger       *zap.Logger
	Interval     time.Duration
}

// Scrape runs all available scraper functions
func (ms *MultiScraper) Scrape() error {
	var lastErr error // Stores the most recent error

	// Loop through each scraper in order
	for _, scraper := range ms.Scrapers {
		// Try scraping the current scraper
		err := scraper.Scrape()

		// If scraping failed, store this error 
		// and continue to the next scraper
		lastErr = err
	}

	return lastErr // Return the last error
}

// Run method begins the periodic scraping in a goroutine.
func (ms *MultiScraper) Run(ctx context.Context) {
	go func() {
		// First run immediately
		for {
			ms.Logger.Info("Scraper cycle started")
			if err := ms.Scrape(); err != nil {
				ms.Logger.Error("Scrape error", zap.Error(err))
			}
			ms.Logger.Info("Scraper cycle finished")

			select {
			case <-time.After(ms.Interval):
				// continue
			case <-ctx.Done():
				fmt.Println("Scraper stopped due to shutdown signal")
				ms.Logger.Info("Scraper stopped due to shutdown signal")
				return
			}
		}
	}()
}
