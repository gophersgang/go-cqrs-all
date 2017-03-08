package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type repoInfo struct {
	url         string
	description string
	lastcommit  string
}

func (ri repoInfo) Markdown() string {
	return fmt.Sprintf("- %s - <br/> %s <br/> ( %s )", ri.url, ri.description, ri.lastcommit)
}

type reposByLastcommit []repoInfo

func (ris reposByLastcommit) Len() int           { return len(ris) }
func (ris reposByLastcommit) Less(i, j int) bool { return ris[i].lastcommit > ris[j].lastcommit }
func (ris reposByLastcommit) Swap(i, j int)      { ris[i], ris[j] = ris[j], ris[i] }

type reposByUrl []repoInfo

func (ris reposByUrl) Len() int           { return len(ris) }
func (ris reposByUrl) Less(i, j int) bool { return ris[i].url < ris[j].url }
func (ris reposByUrl) Swap(i, j int)      { ris[i], ris[j] = ris[j], ris[i] }

func main() {
	var wg sync.WaitGroup
	urls := loadUrls()
	repos := []repoInfo{}
	for _, url := range urls {
		wg.Add(1)
		a := url
		go func() {
			defer wg.Done()
			repo := process(a)
			// fmt.Println(repo)
			repos = append(repos, repo)
			fmt.Print(".")
		}()
	}
	wg.Wait()

	sort.Sort(reposByUrl(repos))
	fmt.Print("\n\n")
	for _, r := range repos {
		fmt.Println(r.Markdown())
	}
}

func process(url string) repoInfo {
	doc := getDoc(url)
	repo := repoInfo{
		url:         strings.ToLower(url),
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
	if hasIncludedLastcommit(doc) {
		return getLastcommitIncluded(doc)
	}
	return getLastcommitAjax(doc)
}

func hasIncludedLastcommit(doc *goquery.Document) bool {
	found := true
	doc.Find(".commit-loader").Each(func(i int, s *goquery.Selection) {
		found = false
	})
	return found
}

func getLastcommitIncluded(doc *goquery.Document) string {
	var datetime string
	doc.Find(".commit-tease relative-time").Each(func(i int, s *goquery.Selection) {
		datetime, _ = s.Attr("datetime")
	})
	return datetime
}

func getLastcommitAjax(doc *goquery.Document) string {
	// extract the ajax url
	// e.g.: <include-fragment class="commit-tease commit-loader" src="/f2prateek/coi/tree-commit/866dee22e2b11dd9780770c00bae53886d9b4863">
	s := doc.Find(".commit-loader")
	path, _ := s.Attr("src")
	url := "https://github.com" + path
	ajaxDoc := urlDoc(url)
	return getLastcommit(ajaxDoc)
}

func getDoc(url string) *goquery.Document {
	return urlDoc(url)
	// return localDoc()
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
