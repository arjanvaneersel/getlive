package scraper

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
	// "github.com/pkg/errors"
	// gourl "net/url"
)

type Page struct {
	mu          sync.RWMutex
	URL         string
	Title       string
	Description string
	Subpages    []Page
	Count       int
}

func (p *Page) addSubPage(page Page) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Subpages = append(p.Subpages, page)
}

type Scraper struct {
	mu    sync.RWMutex
	urls  []string
	Pages []Page
}

func New(urls ...string) *Scraper {
	return &Scraper{
		urls: urls,
	}
}

func (s *Scraper) GetPage(url string, deep bool) (*Page, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	page := Page{
		URL:   url,
		Count: 1,
	}
	page.Title = doc.Find("title").First().Text()
	page.Description = doc.Find("description").First().Text()

	if !deep {
		return &page, nil
	}

	// Get the subpages concurrently for deep searching
	var wg sync.WaitGroup
	doc.Find("a").Each(func(i int, e *goquery.Selection) {
		wg.Add(1)
		go func(p *Page) {
			defer wg.Done()
			href, ok := e.Attr("href")
			if !ok {
				return
			}

			subpage, err := s.GetPage(href, false)
			if err != nil {
				return
			}
			p.addSubPage(*subpage)
		}(&page)
	})
	wg.Wait()

	return &page, nil
}

func (s *Scraper) addPage(p *Page) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Pages = append(s.Pages, *p)
}

func (s *Scraper) Run(pageChan chan *Page) {
	var wg sync.WaitGroup
	fmt.Printf("Starting to scrape %d urls\n", len(s.urls))
	for i, url := range s.urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			fmt.Printf("\t[%d] Scraping %q\n", i, url)
			page, err := s.GetPage(url, true)
			if err != nil {
				fmt.Printf("error getting page %q: %v\n", url, err)
				return
			}

			// If a channel is provided send the result to the channel,
			// else store it in the scraper
			if pageChan != nil {
				fmt.Printf("\t[%d] Sending page %q via channel\n", i, page.Title)
				pageChan <- page
			} else {
				s.addPage(page)
			}
		}(i, url)
	}
	wg.Wait()

	// Close the channel, if provided
	if pageChan != nil {
		close(pageChan)
	}
}

// func (s *Scraper) addURL(url string, wg *sync.WaitGroup) {
// 	if site, ok := s.found[url]; ok {
// 		s.mu.Lock()
// 		site.Count++
// 		s.mu.Unlock()
// 	} else {
// 		// Parse HTML
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			site, err := s.parseHTML(url)
// 			if err != nil {
// 				fmt.Printf("ignoring %s, because of error: %v\n", url, err)
// 				return
// 			}
// 			s.mu.Lock()
// 			s.found[url] = *site
// 			s.mu.Unlock()
// 		}()
// 	}
// }

// func (s *Scraper) parseHTML(url string) (*Site, error) {
// 	uri, err := gourl.Parse(url)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "parsing URL string")
// 		os.Exit(1)
// 	}

// 	fmt.Println("Parsing:", uri.String())
// 	res, err := http.Get(uri.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
// 	}

// 	doc, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	site := Site{
// 		URL:   url,
// 		Count: 1,
// 	}

// 	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
// 		// For each item found, get the band and title
// 		href, ok := sel.Find("a").First().Attr("href")
// 		if ok {
// 			site.addLink(href)
// 		}
// 		// title := s.Find("i").Text()
// 	})

// 	return &site, nil
// }

// func (s *Scraper) worker(i int, wg *sync.WaitGroup, url string) {
// 	defer wg.Done()
// 	fmt.Printf("worker %d starting\n", i)
// 	c := colly.NewCollector()

// 	// Find and visit all links
// 	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
// 		s.addURL(e.Attr("href"), wg)
// 	})

// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting", r.URL)
// 	})

// 	c.Visit(url)
// 	fmt.Printf("worker %d done\n", i)
// }

// func (s *Scraper) Start() (interface{}, error) {
// 	var wg sync.WaitGroup

// 	for i, url := range s.urls {
// 		wg.Add(1)
// 		go s.worker(i, &wg, url)
// 	}

// 	wg.Wait()

// 	for url, v := range s.found {
// 		fmt.Printf("%s was referenced %d times\n", url, v.Count)
// 	}

// 	return nil, nil
// }
