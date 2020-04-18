package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/arjanvaneersel/getlive/internal/scraper"
)

func main() {
	fmt.Sprintf("Scraper")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	scraper.New(
		"https://www.billboard.com/articles/columns/pop/9335531/coronavirus-quarantine-music-events-online-streams",
		"https:www.vulture.com/amp/2020/04/all-musicians-streaming-live-concerts.html",
		"https://www.npr.org/2020/03/17/816504058/a-list-of-live-virtual-concerts-to-watch-during-the-coronavirus-shutdown",
	).Run(nil)

	<-sigChan
	fmt.Printf("shutting down...")
}
