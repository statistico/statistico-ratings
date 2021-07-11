package fixture

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-football-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"time"
)

type Fetcher interface {
	ByCompetition(ctx context.Context, competitionID uint64, numSeasons int8) ([]*statistico.Fixture, error)
	ByDate(ctx context.Context, date time.Time) ([]*statistico.Fixture, error)
}

type fetcher struct {
	competitions   []uint64
	fixtureClient  statisticodata.FixtureClient
	seasonClient   statisticodata.SeasonClient
	clock          clockwork.Clock
}

func (f *fetcher) ByCompetition(ctx context.Context, competitionID uint64, numSeasons int8) ([]*statistico.Fixture, error) {
	res, err := f.seasonClient.ByCompetitionID(ctx, competitionID, "name_desc")

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

func (f *fetcher) ByDate(ctx context.Context, date time.Time) ([]*statistico.Fixture, error) {
	year, month, day := date.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, month, day, 23, 59, 59, 0, time.UTC)

	request := statistico.FixtureSearchRequest{
		DateAfter:            &wrappers.StringValue{Value: start.Format(time.RFC3339)},
		DateBefore:           &wrappers.StringValue{Value: end.Format(time.RFC3339)},
		Sort:                 &wrappers.StringValue{Value: "date_asc"},
	}

	response, err := f.fixtureClient.Search(ctx, &request)

	if err != nil {
		return nil, err
	}

	var fixtures []*statistico.Fixture

	for _, fixture := range response {
		for _, competition := range f.competitions {
			if fixture.Competition.Id == competition {
				fixtures = append(fixtures, fixture)
			}
		}
	}

	return fixtures, nil
}

func parseSeasons(s []*statistico.Season, num int8) []uint64 {
	x := s[len(s)-int(num):]

	var seasons []uint64

	for _, season := range x {
		seasons = append(seasons, season.Id)
	}

	return seasons
}

func NewFetcher(c []uint64, f statisticodata.FixtureClient, s statisticodata.SeasonClient, cl clockwork.Clock) Fetcher {
	return &fetcher{
		competitions:  c,
		fixtureClient: f,
		seasonClient:  s,
		clock:         cl,
	}
}
