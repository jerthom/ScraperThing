package show

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gocolly/colly"

	"github.com/jerthom/ScraperThing/actor"
)

var imdbCollector = colly.NewCollector(colly.AllowedDomains("imdb.com", "www.imdb.com"))

type Show struct {
	Title  string        `json:"title"`
	URL    string        `json:"url"`
	Actors []actor.Actor `json:"actors"`
}

func NewShow(showUrl string) (*Show, error) {
	retS := &Show{URL: showUrl}

	a, err := actors(showUrl)
	// Sort the actors by name; this makes comparing show actor lists quicker
	sort.Slice(a, func(i int, j int) bool {
		return a[i].Name < a[j].Name
	})

	if err != nil {
		return nil, fmt.Errorf("error getting show details: %w", err)
	}
	retS.Actors = a

	return retS, nil
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
					a.URL = e.Request.AbsoluteURL(partialUrl)
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

func SharedActors(s1 *Show, s2 *Show) []actor.Actor {
	var shared []actor.Actor

	for i, j := 0, 0; i < len(s1.Actors) && j < len(s2.Actors); {
		if s1.Actors[i].Name == s2.Actors[j].Name {
			shared = append(shared, s1.Actors[i])
			i += 1
			j += 1
			continue
		}

		if s1.Actors[i].Name < s2.Actors[j].Name {
			i += 1
			continue
		}

		if s1.Actors[i].Name > s2.Actors[j].Name {
			j += 1
			continue
		}
	}

	return shared
}
