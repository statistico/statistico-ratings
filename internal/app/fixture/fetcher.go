package fixture

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-football-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"time"
)

type Fetcher interface {
	ByCompetition(ctx context.Context, competitionID, seasonID uint64) ([]*statistico.Fixture, error)
	ByDate(ctx context.Context, from, to time.Time) ([]*statistico.Fixture, error)
}

type fetcher struct {
	competitions  []uint64
	fixtureClient statisticodata.FixtureClient
	seasonClient  statisticodata.SeasonClient
	clock         clockwork.Clock
}

func (f *fetcher) ByCompetition(ctx context.Context, competitionID, seasonID uint64) ([]*statistico.Fixture, error) {
	res, err := f.seasonClient.ByCompetitionID(ctx, competitionID, "name_desc")

	if err != nil {
		return nil, err
	}

	season, err := parseSeason(res, seasonID)

	if err != nil {
		return nil, err
	}

	req := statistico.FixtureSearchRequest{
		SeasonIds:  []uint64{season.Id},
		DateBefore: &wrappers.StringValue{Value: f.clock.Now().Format(time.RFC3339)},
		Sort:       &wrappers.StringValue{Value: "date_asc"},
	}

	return f.fixtureClient.Search(ctx, &req)
}

func (f *fetcher) ByDate(ctx context.Context, from, to time.Time) ([]*statistico.Fixture, error) {
	request := statistico.FixtureSearchRequest{
		DateAfter:  &wrappers.StringValue{Value: from.Format(time.RFC3339)},
		DateBefore: &wrappers.StringValue{Value: to.Format(time.RFC3339)},
		Sort:       &wrappers.StringValue{Value: "date_asc"},
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

func parseSeason(s []*statistico.Season, id uint64) (*statistico.Season, error) {
	for _, season := range s {
		if season.Id == id {
			return season, nil
		}
	}

	return nil, fmt.Errorf("season %d does not exist", id)
}

func NewFetcher(c []uint64, f statisticodata.FixtureClient, s statisticodata.SeasonClient, cl clockwork.Clock) Fetcher {
	return &fetcher{
		competitions:  c,
		fixtureClient: f,
		seasonClient:  s,
		clock:         cl,
	}
}
