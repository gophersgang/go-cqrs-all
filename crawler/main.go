package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	doc := getDoc()
	// Find description
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name == "description" {
			content, _ := s.Attr("content")
			fmt.Println(content)
		}
	})
}

func getDoc() *goquery.Document {
	// return localDoc()
	return urlDoc()
}

func urlDoc() *goquery.Document {
	url := "https://github.com/andrewwebber/cqrs"
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
