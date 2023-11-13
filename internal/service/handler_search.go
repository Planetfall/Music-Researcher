package service

import (
	"context"
	"fmt"

	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
)

func (s *Service) Search(ctx context.Context, params *pb.Parameters) (*pb.Results, error) {

	trackList, err := s.mySpotify.Search(ctx,
		params.Query,
		params.GenreFilters,
		int(params.Limit),
	)
	if err != nil {
		s.srv.Raise(
			fmt.Sprintf("failed to search spotify with params: %v", params),
			err, nil)
		return nil, err
	}

	return &pb.Results{
		Albums:  nil,
		Artists: nil,
		Tracks:  trackList,
	}, nil
}
