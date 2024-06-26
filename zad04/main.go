package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type TeamStanding struct {
	Position string
	Team     string
	Played   string
	Points   string
}

func main() {
	standings := []TeamStanding{}

	c := colly.NewCollector()

	c.OnHTML("table.wikitable:nth-of-type(1) > tbody", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			position := el.ChildText("td:nth-child(1)")
			team := el.ChildText("td:nth-child(2) a")
			played := el.ChildText("td:nth-child(3)")
			points := el.ChildText("td:nth-child(6)")
			if position != "" && team != "" && played != "" && points != "" {
				standings = append(standings, TeamStanding{
					Position: position,
					Team:     team,
					Played:   played,
					Points:   points,
				})
			}
		})
	})

	c.Visit("https://en.wikipedia.org/wiki/Ekstraklasa")

	file, err := os.Create("ekstraklasa_standings.csv")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Position", "Team", "Seasons", "Points"})

	for _, standing := range standings {
		writer.Write([]string{standing.Position, standing.Team, standing.Played, standing.Points})
	}

	fmt.Println("finished")
}

