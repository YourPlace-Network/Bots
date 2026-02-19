package src

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}
type RSSChannel struct {
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
	Title       string    `xml:"title"`
}
type RSSItem struct {
	Description string `xml:"description"`
	GUID        string `xml:"guid"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Title       string `xml:"title"`
}

var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
var whitespaceRegex = regexp.MustCompile(`\s+`)

func FetchFeed(url string) ([]RSSItem, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch feed %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("feed %s returned status %d", url, resp.StatusCode)
	}
	var feed RSSFeed
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(&feed); err != nil {
		return nil, fmt.Errorf("could not parse feed %s: %w", url, err)
	}
	for i := range feed.Channel.Items {
		feed.Channel.Items[i].Description = sanitizeHTML(feed.Channel.Items[i].Description)
		feed.Channel.Items[i].Title = sanitizeHTML(feed.Channel.Items[i].Title)
		if feed.Channel.Items[i].GUID == "" {
			feed.Channel.Items[i].GUID = feed.Channel.Items[i].Link
		}
	}
	return feed.Channel.Items, nil
}

func sanitizeHTML(s string) string {
	s = htmlTagRegex.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = whitespaceRegex.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}
func sanitizeNonPrintable(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 32 || r == '\n' || r == '\r' || r == '\t' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
