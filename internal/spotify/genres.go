package myspotify

import (
	"context"
	"fmt"

	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
)

func (s *MySpotifyImpl) GetGenreList(ctx context.Context) (*pb.GenreList, error) {

	if err := s.refresh(ctx); err != nil {
		return nil, err
	}

	genreList, err := s.client.GetAvailableGenreSeeds(ctx)
	if err != nil {
		return nil, fmt.Errorf("client.GetAvailableGenreSeeds: %v", genreList)
	}

	return &pb.GenreList{
		Genres: genreList,
	}, nil
}
