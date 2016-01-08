package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	inputFile  = []string{}
	outputFile *os.File
)

func init() {
	inputFile = loadFile("urls")
}

func main() {
	var err error
	outputFile, err = os.Create("output")
	if err != nil {
		fmt.Println(err, "can't create file")
	}
	defer outputFile.Close()

	for _, line := range inputFile {
		if line == "" {
			continue
		}

		fileLine, outLine := updateLine(line)

		_, err := outputFile.WriteString(fileLine)
		if err != nil {
			fmt.Println("write string to file. error: ", err)
		}

		os.Stdout.WriteString(outLine)
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

	// url := doc.Find(".yt-uix-sessionlink").First()
	// fmt.Println("url", url.Find("").Attr("href"))
	// fmt.Printf("%+v\n", url)

	out := ""
	doc.Find(".channels-content-item").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, exist := s.Find(".yt-uix-sessionlink").Attr("href")
		if exist == false {
			fmt.Println("href doesn't exist")
		}
		out = url
		return false
	})

	return out
}

func updateLine(line string) (string, string) {
	ln := strings.Split(line, " ")
	recentURL := getRecentURL(ln[0])
	if len(ln) == 1 || recentURL != ln[1] {
		return ln[0] + " " + recentURL + "\n", ln[0] + "\n"
	}

	return line + "\n", ""
}
