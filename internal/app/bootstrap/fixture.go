package bootstrap

import "github.com/statistico/statistico-ratings/internal/app/fixture"

func (c Container) FixtureFetcher() fixture.Fetcher {
	return fixture.NewFetcher(
		c.Config.SupportedCompetitions,
		c.DataFixtureClient(),
		c.DataSeasonClient(),
		c.Clock,
	)
}
