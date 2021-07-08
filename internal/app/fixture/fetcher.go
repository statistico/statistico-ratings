package fixture

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-football-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/bootstrap"
	"time"
)

type Fetcher struct {
	config         bootstrap.Config
	fixtureClient  statisticodata.FixtureClient
	seasonClient   statisticodata.SeasonClient
	clock          clockwork.Clock
}

func (f *Fetcher) ByCompetition(ctx context.Context, seasonID uint64, numSeasons int8) ([]*statistico.Fixture, error) {
	res, err := f.seasonClient.ByCompetitionID(ctx, seasonID, "name_desc")

	if err != nil {
		return nil, err
	}

	req := statistico.FixtureSearchRequest{
		SeasonIds:            parseSeasons(res, numSeasons),
		DateBefore:           &wrappers.StringValue{Value: f.clock.Now().Format(time.RFC3339)},
		Sort:                 &wrappers.StringValue{Value: "date_asc"},
	}

	return f.fixtureClient.Search(ctx, &req)
}

func (f *Fetcher) ByDate(ctx context.Context, date time.Time) ([]*statistico.Fixture, error) {
	var seasons []uint64

	for _, c := range f.config.SupportedCompetitions {
		res, err := f.seasonClient.ByCompetitionID(ctx, c, "name_desc")

		if err != nil {
			return nil, err
		}

		parsed := parseSeasons(res, 1)

		if len(parsed) == 0 {
			return nil, fmt.Errorf("no seasons returned for competition %d", c)
		}

		seasons = append(seasons, parsed[0])
	}

	year, month, day := date.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, nil)
	end := time.Date(year, month, day, 23, 59, 59, 0, nil)

	req := statistico.FixtureSearchRequest{
		SeasonIds:            seasons,
		DateAfter:            &wrappers.StringValue{Value: start.Format(time.RFC3339)},
		DateBefore:           &wrappers.StringValue{Value: end.Format(time.RFC3339)},
		Sort:                 &wrappers.StringValue{Value: "date_asc"},
	}

	return f.fixtureClient.Search(ctx, &req)
}

func parseSeasons(s []*statistico.Season, num int8) []uint64 {
	x := s[num:]

	var seasons []uint64

	for _, season := range x {
		seasons = append(seasons, season.Id)
	}

	return seasons
}
