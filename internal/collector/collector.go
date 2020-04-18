package collector

import (
	"fmt"
	"time"
)

type Entry struct {
	ID               string    `json:"id"`
	Time             time.Time `json:"datetime"`
	Category         string    `json:"category"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	URL              string    `json:"url"`
	SocialMediaLinks []string  `json:"socialmedia_links"`
	Approved         bool      `json:"approved"`
}

type Collector struct{}

func (c *Collector) Run(done chan struct{}) chan Entry {
	entryChan := make(chan Entry)

	go func() {
		for {
			select {
			case <-done:
				return
			case entry := <-entryChan:
				fmt.Println(entry)
			}
		}
	}()

	return entryChan
}
