package main

import (
	"fmt"

	"github.com/jerthom/ScraperThing/show"
)

func main() {
	//cmd.Execute()
	s1, err := show.NewShow("https://www.imdb.com/title/tt0110912/fullcredits")
	if err != nil {
		fmt.Println(err)
	}

	s2, err := show.NewShow("https://www.imdb.com/title/tt3460252/fullcredits")
	if err != nil {
		fmt.Println(err)
	}

	sharedActors := show.SharedActors([]show.Show{*s1, *s2})

	for _, a := range sharedActors {
		fmt.Println(a)
	}

}
