package scraper

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	mu    sync.RWMutex
	urls  []string
	found map[string]uint
}

func New(urls ...string) *Scraper {
	return &Scraper{
		urls:  urls,
		found: make(map[string]uint),
	}
}

func (b *Scraper) addURL(url string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.found[url]; ok {
		b.found[url]++
	} else {
		b.found[url] = 1
	}
}

func (b *Scraper) worker(i int, wg *sync.WaitGroup, url string) {
	defer wg.Done()
	fmt.Printf("worker %d starting\n", i)
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		b.addURL(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)
	fmt.Printf("worker %d done\n", i)
}

func (b *Scraper) Start() (interface{}, error) {
	var wg sync.WaitGroup

	for i, url := range b.urls {
		wg.Add(1)
		go b.worker(i, &wg, url)
	}

	wg.Wait()

	for url, v := range b.found {
		fmt.Printf("%s was referenced %d times\n", url, v)
	}

	return nil, nil
}
