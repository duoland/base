package rss

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestParseRSSBody(t *testing.T) {
	rssBody, _ := ioutil.ReadFile("./antirez.rss")
	xmlDecoder := xml.NewDecoder(bytes.NewBuffer([]byte(rssBody)))
	var rssObj RSS
	err := xmlDecoder.Decode(&rssObj)
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, item := range rssObj.Channel.Items {
		t.Log(item.Title, item.Link, item.PubDate)
	}
}
