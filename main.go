package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/jerthom/ScraperThing/show"
)

func main() {
	//cmd.Execute()
	start := time.Now()
	showUrls := []string{"https://www.imdb.com/title/tt0110912/fullcredits", "https://www.imdb.com/title/tt3460252/fullcredits"}
	showChan := make(chan show.Show, len(showUrls))
	var wg sync.WaitGroup
	for _, showUrl := range showUrls {
		wg.Add(1)
		go show.NewShow(showUrl, showChan, &wg)
	}
	wg.Wait()
	close(showChan)

	var shows []show.Show
	for s := range showChan {
		shows = append(shows, s)
	}
	sharedActors := show.SharedActors(shows)

	for _, a := range sharedActors {
		fmt.Println(a)
	}
	end := time.Now()
	duration := end.Sub(start)
	fmt.Println("main duration: ", duration)
}
