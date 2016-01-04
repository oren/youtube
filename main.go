package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
	doc, err := goquery.NewDocument("https://www.youtube.com/user/LastWeekTonight/videos")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("doc", doc)

	doc.Find(".channels-content-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".yt-lockup-title").Text()
		url, exist := s.Find(".yt-uix-sessionlink").Attr("href")
		if exist == false {
			fmt.Println("href doesn't exist")
		}
		fmt.Println("title", title)
		fmt.Println("url", url)
		fmt.Println("\n")
	})
}

func main() {
	ExampleScrape()
}
