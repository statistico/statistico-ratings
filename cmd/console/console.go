package main

import (
	"context"
	"fmt"
	"github.com/statistico/statistico-ratings/internal/app/bootstrap"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := bootstrap.BuildContainer(bootstrap.BuildConfig())
	handler := app.TeamRatingHandler()
	ctx := context.Background()

	console := &cli.App{
		Name: "Statistico Ratings - Command Line Application",
		Commands: []cli.Command{
			{
				Name:        "team:by-competition",
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
					handler.ByCompetition(ctx, c.Uint64("competition"), c.Uint64("season"))
					return nil
				},
				Flags: []cli.Flag{
					&cli.Uint64Flag{
						Name:     "competition",
						Usage:    "The ID of the competition to calculate ratings for",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "season",
						Usage:    "The ID of the season to calculate ratings for",
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
