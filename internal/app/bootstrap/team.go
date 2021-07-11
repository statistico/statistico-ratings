package bootstrap

import "github.com/statistico/statistico-ratings/internal/app/team"

func (c Container) TeamRatingCalculator() team.RatingCalculator {
	return team.NewRatingCalculator(c.DataEventClient(), c.Clock)
}

func (c Container) TeamRatingHandler() team.RatingHandler {
	return team.NewHandler(c.FixtureFetcher(), c.TeamRatingProcessor(), c.Logger)
}

func (c Container) TeamRatingProcessor() team.RatingProcessor {
	return team.NewRatingProcessor(
		c.TeamRatingReader(),
		c.TeamRatingWriter(),
		c.TeamRatingCalculator(),
	)
}

func (c Container) TeamRatingReader() team.RatingReader {
	return team.NewRatingReader(c.Database)
}

func (c Container) TeamRatingWriter() team.RatingWriter {
	return team.NewRatingWriter(c.Database)
}
