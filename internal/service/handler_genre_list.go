package service

import (
	"context"

	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
)

func (s *Service) GetGenreList(ctx context.Context, empty *pb.Empty) (*pb.GenreList, error) {

	genreList, err := s.mySpotify.GetGenreList(ctx)
	if err != nil {
		s.srv.Raise(
			"failed to get genre list from spotify",
			err, nil)
		return nil, err
	}

	return genreList, nil
}
