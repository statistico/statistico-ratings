package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/statistico/statistico-ratings/internal/app/bootstrap"
	"github.com/urfave/cli"
	"io"
	"os"
	"strconv"
)

func main() {
	app := bootstrap.BuildContainer(bootstrap.BuildConfig())
	handler := app.TeamRatingHandler()
	ctx := context.Background()

	console := &cli.App{
		Name: "Statistico Ratings - Command Line Application",
		Commands: []cli.Command{
			{
				Name:        "team:csv",
				Usage:       "Calculate team ratings for a competition and season",
				Description: "Calculate team ratings for a competition and season",
				Before: func(c *cli.Context) error {
					fmt.Println("Calculating team ratings...")
					return nil
				},
				After: func(c *cli.Context) error {
					fmt.Println("Complete.")
					return nil
				},
				Action: func(c *cli.Context) error {
					cs, err := os.Open(c.String("filepath"))

					if err != nil {
						return err
					}

					seasons := csv.NewReader(cs)

					for {
						row, err := seasons.Read()

						if err == io.EOF {
							break
						}

						comp, _ := strconv.ParseUint(row[0], 0, 64)
						season, _ := strconv.ParseUint(row[1], 0, 64)

						handler.ByCompetition(ctx, comp, season)
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "filepath",
						Usage:    "The filepath where the csv resides",
						Required: true,
					},
				},
			},
			{
				Name:        "team:today",
				Usage:       "Calculate team ratings for today's fixtures",
				Description: "Calculate team ratings for today's fixtures",
				Before: func(c *cli.Context) error {
					fmt.Println("Calculating team ratings...")
					return nil
				},
				After: func(c *cli.Context) error {
					fmt.Println("Complete.")
					return nil
				},
				Action: func(c *cli.Context) error {
					err := handler.Today(ctx, c.Int("hour"))

					if err != nil {
						fmt.Printf("error: %s\n", err.Error())
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "hour",
						Usage:    "Process the days fixture played before the given hour",
						Required: true,
					},
				},
			},
		},
	}

	err := console.Run(os.Args)

	if err != nil {
		fmt.Printf("Error in executing command: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
