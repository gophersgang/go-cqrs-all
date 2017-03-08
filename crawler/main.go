package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type repoInfo struct {
	description string
	lastcommit  string
}

func main() {
	var wg sync.WaitGroup
	urls := loadUrls()
	for _, url := range urls {
		wg.Add(1)
		a := url
		go func() {
			defer wg.Done()
			desc := process(a)
			fmt.Println(desc)
		}()
	}
	wg.Wait()
}

func process(url string) repoInfo {
	doc := getDoc(url)
	repo := repoInfo{
		description: getDescription(doc),
		lastcommit:  getLastcommit(doc),
	}
	return repo
}

func getDescription(doc *goquery.Document) string {
	var content string
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name == "description" {
			content, _ = s.Attr("content")
		}
	})
	return content
}

func getLastcommit(doc *goquery.Document) string {
	var datetime string
	doc.Find(".commit-tease relative-time").Each(func(i int, s *goquery.Selection) {
		datetime, _ = s.Attr("datetime")
	})
	return datetime
}

func getDoc(url string) *goquery.Document {
	return urlDoc(url)
}

func urlDoc(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	check(err)
	return doc
}

func localDoc() *goquery.Document {
	filename := "crawler/fixture.html"
	file, err := os.Open(filename)
	check(err)
	doc, err := goquery.NewDocumentFromReader(file)
	check(err)
	return doc
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadUrls() []string {
	return file2lines("data/urls.txt")
}

/*
simple lines reader
*/
func file2lines(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if t := scanner.Text(); validURL(t) {
			lines = append(lines, t)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return lines
}

func validURL(l string) bool {
	return !strings.Contains(l, " ") && len(l) != 0
}
