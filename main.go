package main

import (
	"fmt"

	"github.com/jerthom/ScraperThing/show"
)

func main() {
	actors, err := show.ActorsForShow("https://www.imdb.com/title/tt11737520/fullcredits")
	if err != nil {
		fmt.Println(err)
	}

	for _, a := range actors {
		fmt.Println(a)
	}
}
