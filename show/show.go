package show

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"

	"github.com/jerthom/ScraperThing/actor"
)

type Show struct {
	Title  string        `json:"title"`
	URL    string        `json:"url"`
	Actors []actor.Actor `json:"actors"`
}

func ActorsForShow(showUrl string) ([]actor.Actor, error) {
	var actorURLs []actor.Actor

	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

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

//func SharedActors(s1 Show, s2 Show) []actor.Actor {
//	var shared []actor.Actor
//
//	for a := range s1.Actors {
//
//	}
//}
