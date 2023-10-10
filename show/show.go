package show

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"

	"github.com/jerthom/ScraperThing/actor"
)

var imdbCollector = colly.NewCollector(colly.AllowedDomains("imdb.com", "www.imdb.com"))

type Show struct {
	Title  string        `json:"title"`
	URL    string        `json:"url"`
	Actors []actor.Actor `json:"actors"`
}

// NewShowCon (concurrent) creates a Show by scraping the relevant data from the provided url, and then writes that show to showChan
func NewShowCon(showUrl string, showChan chan Show, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()

	// Scrape actor information
	a, err := actors(showUrl)
	if err != nil {
		fmt.Println("error getting show details: ", err)
		return
	}

	s := Show{URL: showUrl, Actors: a}
	showChan <- s
	end := time.Now()
	diff := end.Sub(start)
	fmt.Println("duration ", showUrl, ": ", diff)
}

// NewShowSec (sequential) creates a Show by scraping the relevant data from the provided url
func NewShowSec(showUrl string) (*Show, error) {
	start := time.Now()

	// Scrape actor information
	a, err := actors(showUrl)
	if err != nil {
		return nil, fmt.Errorf("error getting show details: %w", err)
	}

	s := Show{URL: showUrl, Actors: a}
	end := time.Now()
	diff := end.Sub(start)
	fmt.Println("duration ", showUrl, ": ", diff)
	return &s, nil
}

func actors(showUrl string) ([]actor.Actor, error) {
	var actorURLs []actor.Actor

	c := imdbCollector.Clone()

	c.OnHTML("table.cast_list > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			el.ForEach("td", func(k int, el2 *colly.HTMLElement) {
				if k == 1 {
					a := actor.Actor{}
					a.Name = strings.TrimSpace(el2.Text)

					partialUrl := el2.ChildAttr("a[href]", "href")
					url := e.Request.AbsoluteURL(partialUrl)
					// Trim reference off end of url.
					url, _, _ = strings.Cut(url, "/?ref")
					a.URL = url

					actorURLs = append(actorURLs, a)
				}
			})
		})
	})

	err := c.Visit(showUrl)
	if err != nil {
		return nil, fmt.Errorf("error getting actor urls: %w", err)
	}
	return actorURLs, nil
}

// SharedActors returns the list of Actors which appear in *at least* the number of shows indicated by threshold
func SharedActors(shows []Show, threshold int) []actor.Actor {
	actors := make(map[actor.Actor]int)
	for _, s := range shows {
		for _, a := range s.Actors {
			actors[a] = actors[a] + 1
		}
	}
	var shared []actor.Actor
	for a, i := range actors {
		if i >= threshold {
			shared = append(shared, a)
		}
	}
	return shared
}
