package rss

import (
	"encoding/xml"
	"time"
)

// RSS is the root node of the rss body
type RSS struct {
	xml.Name `xml:"rss"`
	Version  string  `xml:"version,attr"`
	Channel  Channel `xml:"channel"`
}

// Channel is the root node of the rss channel
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item is the root node of the rss item
type Item struct {
	Title       string `xml:"title"`
	GUID        string `xml:"guid"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"` // time.RFC1123Z
	Description string `xml:"description"`
	Comments    string `xml:"comments"`
}

// ParsePubDate parse the publish date
func ParsePubDate(pubDate string) (t time.Time) {
	t, _ = time.Parse(time.RFC1123Z, pubDate)
	return
}
