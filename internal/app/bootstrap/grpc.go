package bootstrap

import "github.com/statistico/statistico-ratings/internal/app/grpc"

func (c Container) GrpcTeamRatingService() *grpc.TeamRatingService {
	return grpc.NewTeamRatingService(c.TeamRatingReader(), c.Logger)
}
