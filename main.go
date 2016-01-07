package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func loadFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err, "can't read file")
	}
	lines := strings.Split(string(content), "\n")

	return lines
}

func getRecentURL(webpage string) string {
	doc, err := goquery.NewDocument(webpage)
	if err != nil {
		log.Fatal(err)
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
		return ln[0] + " " + recentURL + "\n", ln[0] + " " + recentURL + "\n"
	}

	return line + "\n", ""
}

func main() {
	in := []string{}
	outFile := []string{}
	out := []string{}

	in = loadFile("urls")
	for _, line := range in {
		if line == "" {
			continue
		}

		fileLine, outLine := updateLine(line)

		outFile = append(outFile, fileLine)

		if outLine != "" {
			out = append(out, "youtube.com"+strings.Split(outLine, " ")[1])
		}
	}

	f, err := os.Create("output")
	if err != nil {
		fmt.Println(err, "can't create file")
	}
	defer f.Close()

	for _, s := range outFile {
		_, err := f.WriteString(s)
		if err != nil {
			fmt.Println(err, "write string to file")
		}
	}

	for _, s := range out {
		os.Stdout.WriteString(s)
	}
}
