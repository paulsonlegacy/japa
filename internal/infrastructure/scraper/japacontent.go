package scraper

import (
	"fmt"
	"time"
	
	"japa/internal/app/http/dto/request"
	"japa/internal/domain/entity"

	"github.com/gocolly/colly"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)


var (
	scrapeURL     string = "https://japacontent.com/"
	allowedDomain string = "japacontent.com"
)


type JapaContentScraper struct {
	Logger  *zap.Logger
	DB     *gorm.DB
}


func NewJapaContentScraper(logger *zap.Logger, db *gorm.DB) *JapaContentScraper {
	return &JapaContentScraper{
		DB: db,
		Logger: logger,
	}
}


func (s *JapaContentScraper) Scrape() error {
	var lastErr error // Stores the most recent error

	scrapedItems, err := s.ScrapeData()
	if err != nil {
		lastErr = err
	}

	requestItems, err := s.ConvertToScrapedPostRequests(scrapedItems)
	if err != nil {
		lastErr = err
	}

	if err := s.SaveScrapedPosts(requestItems); err != nil {
		lastErr = err
	}

	if lastErr == nil {
		s.Logger.Info(
			fmt.Sprintf("Scraping + saving for %s completed", scrapeURL),
			zap.String("url", scrapeURL),
			zap.Time("timestamp", time.Now()),
		)
	}

	return lastErr
}


// Method to scrape the actual data
func (s *JapaContentScraper) ScrapeData() ([]map[string]string, error) {
	// List of scraped items
	var scrapedItems  []map[string]string

	// One collector for the list page
	listCollector := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
	)

	// One collector for article pages
	postCollector := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
	)

	// Log requests
	listCollector.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting list:", r.URL)
		s.Logger.Info(
			"Visiting blog",
			zap.String("url", r.URL.String()),
			zap.Time("timestamp", time.Now()),
		)
	})
	postCollector.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting article:", r.URL)
		s.Logger.Info(
			"Visiting post",
			zap.String("url", r.URL.String()),
			zap.Time("timestamp", time.Now()),
		)
	})

	// Handle errors
	listCollector.OnError(func(r *colly.Response, err error) {
		s.Logger.Error(
			"Scrape blog request error",
			zap.String("url", r.Request.URL.String()),
			zap.Time("timestamp", time.Now()),
		)
	})
	postCollector.OnError(func(r *colly.Response, err error) {
		s.Logger.Error(
			"Scrape post request error",
			zap.String("url", r.Request.URL.String()),
			zap.Time("timestamp", time.Now()),
		)
	})

	// For each post summary in list page
	listCollector.OnHTML("div.td-cpt-post", func(e *colly.HTMLElement) {
		// Build partial result
		post := map[string]string{
			"title":    e.ChildText("h3.entry-title a"),
			"source":  e.ChildAttr("h3.entry-title a", "href"),
			"category": e.ChildText("a.td-post-category"),
		}
		style := e.ChildAttr("span.entry-thumb", "style")
		post["post_img"] = extractBackgroundImageURL(style)

		// Save partial result now
		scrapedItems = append(scrapedItems, post)

		// Visit the post to get more details
		postCollector.Visit(post["source"])
	})

	// In the article page, extract extra fields
	postCollector.OnHTML("div.td-post-content", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		scrapedContentHTML := e.DOM.Clone()

		// Remove known unwanted elements
		scrapedContentHTML.Find("style, script, iframe, .ads, .related-posts").Remove()

		// Get inner HTML
		innerScrapedContentHTML, _ := scrapedContentHTML.Html()
		safeContentText := scrapedContentHTML.Text()


		// Sanitize
		p := bluemonday.UGCPolicy()
		safeContentHTML := p.Sanitize(innerScrapedContentHTML) // Now safeHTML is clean to store or render


		// Try to get meta from head if needed
		mainPostImg, _ := e.DOM.ParentsFiltered("html").Find("meta[property='og:image']").Attr("content")
		description, _ := e.DOM.ParentsFiltered("html").Find("meta[name='description']").Attr("content")
		datePublished, _ := e.DOM.ParentsFiltered("html").Find("meta[property='article:published_time']").Attr("content")


		// Find the corresponding result map by URL
		for index, resultItem := range scrapedItems {
			if resultItem["source"] == url {
				scrapedItems[index]["excerpt"] = description
				scrapedItems[index]["published_at"] = datePublished
				scrapedItems[index]["content_html"] = safeContentHTML
				scrapedItems[index]["content_text"] = safeContentText

				if scrapedItems[index]["post_img"] == "" {
					scrapedItems[index]["post_img"] = mainPostImg
				}

				break
			}
		}
	})

	// After all scraping
	postCollector.OnScraped(func(r *colly.Response) {
		//fmt.Println("Finished scraping article: ", r.Request.URL)
		s.Logger.Info(
			"Finished scraping article",
			zap.String("url", r.Request.URL.String()),
			zap.Time("timestamp", time.Now()),
		)
	})

	listCollector.OnScraped(func(r *colly.Response) {
		//fmt.Println("Finished scraping list page.")
		s.Logger.Info(
			"Finished scraping blog",
			zap.String("url", r.Request.URL.String()),
			zap.Time("timestamp", time.Now()),
		)
	})

	// Start scraping
	if err := listCollector.Visit("https://japacontent.com/"); err != nil {
		s.Logger.Error(
			"Error encountered while scraping",
			zap.String("error", err.Error()),
			zap.Time("timestamp", time.Now()),
		)
		return nil, err
	}

	// Wait for all jobs to finish 
	// before exiting function
	listCollector.Wait()
	postCollector.Wait()

	
	s.Logger.Info(
		fmt.Sprintf("Scraping for %s completed", scrapeURL),
		zap.String("url", scrapeURL),
		zap.Time("timestamp", time.Now()),
	)

	/*
	// Print all collected data
	for i, m := range scrapedItems {
		fmt.Printf("Post #%d:\n", i+1)
		for k, v := range m {
			fmt.Printf("  %s: %s\n", k, v)
		}
		fmt.Println("---")
	}
	*/

	return scrapedItems, nil
}


func (s *JapaContentScraper) ConvertToScrapedPostRequests(results []map[string]string) ([]*request.ScrapedPostRequest, error) {
	var requests []*request.ScrapedPostRequest

	for _, r := range results {
		req := &request.ScrapedPostRequest{
			Category:    r["category"],
			Title:       r["title"],
			Excerpt:     strPtrOrNil(r["excerpt"]),
			PostImg:     strPtrOrNil(r["post_img"]),
			ContentHTML: r["content_html"],
			ContentText: strPtrOrNil(r["content_text"]),
			Source:      r["source"],
		}
		requests = append(requests, req)
	}
	return requests, nil
}


func (s *JapaContentScraper) SaveScrapedPosts(posts []*request.ScrapedPostRequest) error {
	if len(posts) == 0 {
		s.Logger.Info("No scraped posts to save")
		return nil
	}

	// 1) Extract all source URLs
	var sources []string
	for _, post := range posts {
		sources = append(sources, post.Source)
	}

	// 2) Query existing sources to avoid duplicates
	var existingSources []string
	if err := s.DB.
		Model(&entity.ScrapedPost{}).
		Where("source IN ?", sources).
		Pluck("source", &existingSources).
		Error; err != nil {
		return fmt.Errorf("error querying existing sources: %w", err)
	}

	existingSet := make(map[string]struct{}, len(existingSources))
	for _, src := range existingSources {
		existingSet[src] = struct{}{}
	}

	// 3) Prepare batch insert slice
	var toInsert []entity.ScrapedPost
	for _, post := range posts {
		if _, exists := existingSet[post.Source]; exists {
			s.Logger.Info("Skipping duplicate post", zap.String("source", post.Source))
			continue
		}
		toInsert = append(toInsert, entity.ScrapedPost{
			Category:    post.Category,
			Title:       post.Title,
			ContentHTML: post.ContentHTML,
			ContentText: post.ContentText,
			Excerpt:     post.Excerpt,
			PostImg:     post.PostImg,
			Source:      post.Source,
			Status:      "pending",
		})
	}

	// 4) Save to DB
	if len(toInsert) > 0 {
		if err := s.DB.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&toInsert).Error; err != nil {
			return fmt.Errorf("error inserting new scraped posts: %w", err)
		}
		s.Logger.Info("Saved new scraped posts", zap.Int("count", len(toInsert)))
	} else {
		s.Logger.Info("No new posts to insert after deduplication")
	}

	return nil
}

/*
After scraping + saving:

- Admin dashboard queries scraped_posts where status = 'pending'

- Admin can preview, edit, or publish

- When publishing, you update status to published
*/