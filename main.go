package main

import (
	"io/ioutil"
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
	lines = loadFile("urls")
	total = len(lines)
	results = make(chan result)
}

func worker(line string, results chan<- result) {
	ln := strings.Split(line, " ")
	recentURL := getRecentURL(ln[0])
	if len(ln) == 1 || recentURL != ln[1] {
		r := result{
			Stdout: ln[0] + "\n",
			Line:   ln[0] + " " + recentURL + "\n",
		}
		results <- r
		return
	}

	r := result{
		Stdout: "",
		Line:   line + "\n",
	}

	results <- r
}

func main() {
	var err error
	outputFile, err = os.Create("output")
	if err != nil {
		log.Println(err, "can't create file")
	}
	defer outputFile.Close()

	for _, line := range lines {
		if line == "" {
			continue
		}

		go worker(line, results)
	}

	// collect all the results of the workers
	for i := 0; i < total-1; i++ {
		r := <-results
		os.Stdout.WriteString(r.Stdout)
		_, err := outputFile.WriteString(r.Line)
		if err != nil {
			log.Println("write string to file. error: ", err)
		}
	}
}

func loadFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("can't read file. error: ", err)
	}
	lines := strings.Split(string(content), "\n")

	return lines
}

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
