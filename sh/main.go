package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/ghodss/yaml"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Config for urls
type Config struct {
	Urls []string `json:"urls"`
}

func main() {
	fmt.Println(os.Getwd())
	s := newSemaphore(5)
	var wg sync.WaitGroup
	urls := loadUrls()
	wg.Add(len(urls))
	for _, url := range urls {
		s.Acquire(1)
		a := url
		go func() {
			defer wg.Done()
			checkRepo(a)
			defer s.Release(1)
		}()
	}
	wg.Wait()
}

func loadUrls() []string {
	dat, err := ioutil.ReadFile("sh/urls.yml")
	check(err)
	c := Config{}
	yaml.Unmarshal(dat, &c)
	return c.Urls
}

func checkRepo(url string) error {
	repo := newRepo(url)
	repo.Run()
	// time.Sleep(500 * time.Millisecond)
	return nil
}

/*****************************************************************

Repo logic

******************************************************************/

type repo struct {
	url string
}

func newRepo(url string) *repo {
	return &repo{url: url}
}

// initial git checkout
func (r *repo) checkout() error {
	fmt.Printf("checking out %s\n", r.fullPath())
	cmd := fmt.Sprintf("git clone %s %s", r.url, r.fullPath())
	out, err := exec.Command("sh", "-c", cmd).Output()
	check(err)
	fmt.Println(out)
	return nil
}

// refresh existing repo
func (r *repo) refresh() error {
	fmt.Printf("refreshing %s\n", r.fullPath())
	cmd := fmt.Sprintf("cd %s; git pull", r.fullPath())
	_, err := exec.Command("sh", "-c", cmd).Output()
	check(err)
	// fmt.Println(out)
	return nil
}

// does this project exist?
func (r *repo) exists() bool {
	if _, err := os.Stat(r.fullPath()); err == nil {
		return true
	}
	return false
}

func (r *repo) Run() error {
	if r.exists() {
		return r.refresh()
	}
	return r.checkout()
}

// the name of the resulting folder (unique)
func (r *repo) projectName() string {
	parts := strings.Split(r.url, "/")
	user := parts[len(parts)-2]
	name := parts[len(parts)-1]
	name = strings.Replace(name, ".git", "", -1)
	res := fmt.Sprintf("%s--%s", user, name)
	return res
}

// full path to repo folder
func (r *repo) fullPath() string {
	a := []string{"src", r.projectName()}
	return strings.Join(a, "/")
}

/*****************************************************************

Semaphore provides a semaphore synchronization primitive
(vendored for simplicity)

******************************************************************/

// Semaphore controls access to a finite number of resources.
type Semaphore chan struct{}

// New creates a Semaphore that controls access to `n` resources.
func newSemaphore(n int) Semaphore {
	return Semaphore(make(chan struct{}, n))
}

// Acquire `n` resources.
func (s Semaphore) Acquire(n int) {
	for i := 0; i < n; i++ {
		s <- struct{}{}
	}
}

// Release `n` resources.
func (s Semaphore) Release(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}
