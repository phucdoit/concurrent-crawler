package main

import (
	"testing"
)

func TestCrawAndWrite(t *testing.T) {
	type args struct {
		urls []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 1",
			args: args{
				urls: []string{
					"https://www.google.com",
					"https://www.facebook.com",
					"https://www.youtube.com",
					"https://www.wikipedia.org",
					"https://www.amazon.com",
					"https://www.twitter.com",
					"https://www.instagram.com",
					"https://www.linkedin.com",
					"https://www.reddit.com",
					"https://www.pinterest.com",
					"https://www.netflix.com",
					"https://www.microsoft.com",
					"https://www.apple.com",
					"https://www.tumblr.com",
					"https://www.blogger.com",
					"https://www.flickr.com",
					"https://www.yahoo.com",
					"https://www.baidu.com",
					"https://www.yelp.com",
					"https://www.quora.com",
					"https://www.github.com",
					"https://www.stackoverflow.com",
					"https://www.medium.com",
					"https://www.paypal.com",
					"https://www.salesforce.com",
					"https://www.spotify.com",
					"https://www.slack.com",
					"https://www.dropbox.com",
					"https://www.adobe.com",
					"https://www.shopify.com",
					"https://www.twitch.tv",
					"https://www.zoom.us",
					"https://www.airbnb.com",
					"https://www.reuters.com",
					"https://www.cnn.com",
					"https://www.bbc.com",
					"https://www.nytimes.com",
					"https://www.theguardian.com",
					"https://www.washingtonpost.com",
					"https://www.forbes.com",
					"https://www.bloomberg.com",
					"https://www.wsj.com",
					"https://www.ft.com",
					"https://www.economist.com",
					"https://www.vice.com",
					"https://www.buzzfeed.com",
					"https://www.huffpost.com",
					"https://www.techcrunch.com",
					"https://www.engadget.com",
					"https://www.wired.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CrawAndWrite(tt.args.urls)
		})
	}
}
