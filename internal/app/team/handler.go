package team

import (
	"context"
	"github.com/statistico/statistico-football-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/calculate"
	"time"
)

type FixtureHandler interface {
	Handle(ctx context.Context, f *statistico.Fixture, home, away *Rating) (*Rating, *Rating, error)
}

type fixtureHandler struct {
	result statisticodata.ResultClient
	event statisticodata.EventClient
}

func (h *fixtureHandler) Handle(ctx context.Context, f *statistico.Fixture, home, away *Rating) (*Rating, *Rating, error) {
	res, err := h.result.ByID(ctx, uint64(f.Id))

	if err != nil {
		return err
	}

	events, err := h.event.FixtureEvents(ctx, uint64(f.Id))

	if err != nil {
		return err
	}

	hv := calculate.PointsValue(home.Attack.Total, away.Defence.Total, 20, 0)
	av := calculate.PointsValue(away.Attack.Total, home.Defence.Total, 20, 0)
	//
	//home := Rating{
	//	TeamID:    f.HomeTeam.Id,
	//	FixtureID: uint64(f.Id),
	//	SeasonID:  f.Season.Id,
	//	Attack:    Points{
	//		Total:      home.Attack.Total += hv,
	//		Difference: hv,
	//	},
	//	Defence:   Points{},
	//	Timestamp: time.Time{},
	//}

	return nil
}
