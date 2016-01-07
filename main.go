package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// load text file into slice of webpages
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

func loadFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err, "can't read file")
	}
	lines := strings.Split(string(content), "\n")

	return lines
}

func updateLine(line string) string {
	ln := strings.Split(line, " ")
	fmt.Println(ln[0])
	return "wow"
}

func main() {
	in := []string{}
	out := []string{}
	in = loadFile("urls")
	for _, line := range in {
		updatedLine := updateLine(line)
		out = append(out, updatedLine)
	}

	fmt.Println(out)

	// ExampleScrape()
}
