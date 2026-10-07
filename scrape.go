package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// Story holds data scraped from a Hacker News front page row.
type Story struct {
	Rank  string `json:"rank"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Site  string `json:"site"`
}

func main() {
	c := colly.NewCollector()

	var stories []Story

	c.OnHTML("tr.athing", func(e *colly.HTMLElement) {
		stories = append(stories, Story{
			Rank:  strings.TrimSuffix(e.ChildText("span.rank"), "."),
			Title: e.ChildText("span.titleline > a"),
			URL:   e.ChildAttr("span.titleline > a", "href"),
			Site:  e.ChildText("span.sitestr"),
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	if err := c.Visit("https://news.ycombinator.com/"); err != nil {
		fmt.Fprintln(os.Stderr, "visit failed:", err)
		os.Exit(1)
	}

	for _, s := range stories {
		fmt.Printf("%2s. %s (%s)\n    %s\n", s.Rank, s.Title, s.Site, s.URL)
	}

	data, err := json.MarshalIndent(stories, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "marshal failed:", err)
		os.Exit(1)
	}
	fmt.Println(string(data))

	if err := os.WriteFile("output.json", data, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "write failed:", err)
		os.Exit(1)
	}
}
