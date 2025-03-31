package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Define the HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	// Set the request header to custom values
	req.Header.Set("User-Agent", "gator")
	//Create an HTTP client
	client := &http.Client{}
	//Get the HTTP client to perform the request and capture the response
	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()

	//Read the information in the body of the response & store it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	// Create an empty RSSFeed and unmarshal the XML data into it
	feed := &RSSFeed{}

	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return feed, err
	}
	// Remove the escape characters from the feed
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Remove the escape characters from the individual feed items
	for _, feedItem := range feed.Channel.Item {
		feedItem.Title = html.UnescapeString(feedItem.Title)
		feedItem.Description = html.UnescapeString(feedItem.Description)
	}

	return feed, nil
}
