package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/jerthom/ScraperThing/show"
)

func main() {
	start := time.Now()
	//Shows to scrape
	showUrls := []string{"https://www.imdb.com/title/tt0110912/fullcredits", "https://www.imdb.com/title/tt3460252/fullcredits",
		"https://www.imdb.com/title/tt11737520/fullcredits", "https://www.imdb.com/title/tt7131622/fullcredits", "https://www.imdb.com/title/tt0116483/fullcredits",
		"https://www.imdb.com/title/tt0142342/fullcredits"}

	// Scrape for shows concurrently
	showChan := make(chan show.Show, len(showUrls))
	var wg sync.WaitGroup
	for _, showUrl := range showUrls {
		wg.Add(1)
		go show.NewShow(showUrl, showChan, &wg)
	}
	wg.Wait()
	close(showChan)

	// Extract shows from channel
	var shows []show.Show
	for s := range showChan {
		shows = append(shows, s)
	}

	// Get actors who are in at least n shows
	sharedActors := show.SharedActors(shows, 2)

	for _, a := range sharedActors {
		fmt.Println(a)
	}
	end := time.Now()
	duration := end.Sub(start)
	fmt.Println("main duration: ", duration)
}
