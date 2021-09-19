package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

// XXX: refactor these once we get a better sense of the type of crawling we'll
// do and the common things we'll do (e.g., I'm guessing searching for links
// based on regex might be a thing).

func ScrapeLinks(url string) ([]string, error) {
	doc, err := LoadDocument(url)
	if err != nil {
		return nil, err
	}
	links := []string{}
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if !exists {
			// odd.. since this is what we selected...
			return
		}
		links = append(links, href)
	})
	return links, nil
}

func LoadDocument(url string) (*goquery.Document, error) {
	// Grab the HTML
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected status code: %d %s", res.StatusCode, res.Status)
	}
	// and parse it
	return goquery.NewDocumentFromReader(res.Body)
}
