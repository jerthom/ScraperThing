package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	c.OnHTML("table.cast_list > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			el.ForEach("td.primary_photo", func(_ int, el2 *colly.HTMLElement) {
				actorUrl := el2.ChildAttr("a[href]", "href")
				c.Visit(e.Request.AbsoluteURL(actorUrl))
			})
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.imdb.com/title/tt11737520/fullcredits")
}
