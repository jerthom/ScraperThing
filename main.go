package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/jerthom/ScraperThing/show"
)

// Scrape show information concurrently
func showsCon(showUrls []string) []show.Show {
	// Scrape for shows concurrently
	showChan := make(chan show.Show, len(showUrls))
	var wg sync.WaitGroup
	for _, showUrl := range showUrls {
		wg.Add(1)
		go show.NewShowCon(showUrl, showChan, &wg)
	}
	wg.Wait()
	close(showChan)

	// Extract shows from channel
	var shows []show.Show
	for s := range showChan {
		shows = append(shows, s)
	}
	return shows
}

// Scrape show information sequentially
func showsSeq(showUrls []string) []show.Show {
	var shows []show.Show
	for _, sUrl := range showUrls {
		s, err := show.NewShowSec(sUrl)
		if err != nil {
			fmt.Println("Error scraping show ", sUrl, " err: ", err)
			continue
		}
		shows = append(shows, *s)
	}
	return shows
}

func main() {
	start := time.Now()
	//Shows to scrape
	showUrls := []string{"https://www.imdb.com/title/tt0110912/fullcredits", "https://www.imdb.com/title/tt3460252/fullcredits",
		"https://www.imdb.com/title/tt11737520/fullcredits", "https://www.imdb.com/title/tt7131622/fullcredits", "https://www.imdb.com/title/tt0116483/fullcredits",
		"https://www.imdb.com/title/tt0142342/fullcredits", "blarg"}

	parallel := true
	var shows []show.Show
	if parallel {
		shows = showsCon(showUrls)
	} else {
		shows = showsSeq(showUrls)
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
