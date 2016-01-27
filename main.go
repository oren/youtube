package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	lines      = []string{}
	outputFile *os.File
	results    chan result
	total      int
)

type result struct {
	Stdout string
	Line   string
}

func init() {
	results = make(chan result)
}

func main() {
	// 1. setup
	var err error
	outputFile, err = os.Create("output")
	if err != nil {
		log.Println(err, "can't create file")
	}
	defer outputFile.Close()

	inputFile, err := os.Open("urls")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// 2. work
	total := 0
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		total += 1
		go worker(scanner.Text(), results)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// 3. collect
	// collect all the results of the workers
	for i := 0; i < total; i++ {
		r := <-results
		os.Stdout.WriteString(r.Stdout)
		_, err := outputFile.WriteString(r.Line)
		if err != nil {
			log.Println("write string to file. error: ", err)
		}
	}
}

func worker(line string, results chan<- result) {
	ln := strings.Split(line, " ")
	recentURL := getRecentURL(ln[0])
	// if there is a video i didn't watch
	if len(ln) == 1 || recentURL != ln[1] {
		r := result{
			Stdout: ln[0] + "\n",
			Line:   ln[0] + " " + recentURL + "\n",
		}
		results <- r
		return
	}

	// i already watched all the videos of this webpage
	r := result{
		Stdout: "",
		Line:   line + "\n",
	}

	results <- r
}

// fetch webpage and return url for recent video
func getRecentURL(webpage string) string {
	doc, err := goquery.NewDocument(webpage)
	if err != nil {
		log.Fatal("Error scraping webpage: ", err)
	}

	out := ""
	doc.Find(".channels-content-item").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, exist := s.Find(".yt-uix-sessionlink").Attr("href")
		if exist == false {
			log.Println("href doesn't exist")
		}
		out = url
		return false
	})

	return out
}
