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

// NewShow creates a Show by scraping the relevant data from the provided url
func NewShow(showUrl string, shows chan Show, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	retS := Show{URL: showUrl}

	a, err := actors(showUrl)
	if err != nil {
		fmt.Println("error getting show details: %w", err)
	}

	retS.Actors = a
	shows <- retS
	end := time.Now()
	diff := end.Sub(start)
	fmt.Println("duration ", showUrl, ": ", diff)
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
					// Trim reference off end of url.
					url := e.Request.AbsoluteURL(partialUrl)
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

// SharedActors returns the list of Actors which appear in both of the specified shows
func SharedActors(shows []Show) []actor.Actor {
	actors := make(map[actor.Actor]int)
	for _, s := range shows {
		for _, a := range s.Actors {
			actors[a] = actors[a] + 1
		}
	}
	var shared []actor.Actor
	for a, i := range actors {
		if i > 1 {
			shared = append(shared, a)
		}
	}
	return shared
}
