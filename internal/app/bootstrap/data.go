package bootstrap

import (
	"github.com/statistico/statistico-football-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc"
)

func (c Container) DataEventClient() statisticodata.EventClient {
	config := c.Config

	address := config.StatisticoDataService.Host + ":" + config.StatisticoDataService.Port

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		c.Logger.Warnf("Error initializing statistico data service grpc client %s", err.Error())
	}

	client := statistico.NewEventServiceClient(conn)

	return statisticodata.NewEventClient(client)
}

func (c Container) DataFixtureClient() statisticodata.FixtureClient {
	config := c.Config

	address := config.StatisticoDataService.Host + ":" + config.StatisticoDataService.Port

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		c.Logger.Warnf("Error initializing statistico data service grpc client %s", err.Error())
	}

	client := statistico.NewFixtureServiceClient(conn)

	return statisticodata.NewFixtureClient(client)
}

func (c Container) DataSeasonClient() statisticodata.SeasonClient {
	config := c.Config

	address := config.StatisticoDataService.Host + ":" + config.StatisticoDataService.Port

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		c.Logger.Warnf("Error initializing statistico data service grpc client %s", err.Error())
	}

	client := statistico.NewSeasonServiceClient(conn)

	return statisticodata.NewSeasonClient(client)
}
