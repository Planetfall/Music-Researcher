package server

import (
	"context"

	pb "github.com/Dadard29/planetfall/musicresearcher/pkg/pb"
)

func (s *Server) GetGenreList(ctx context.Context, empty *pb.Empty) (*pb.GenreList, error) {
	err := s.setSpotifyClient(ctx)
	if err != nil {
		s.errorReport(err, "failed setting up connection with Spotify")
		return nil, err
	}

	genreList, err := s.spotifyClient.GetAvailableGenreSeeds(ctx)
	if err != nil {
		s.errorReport(err, "failed to retrieve genre list")
		return nil, err
	}
	return &pb.GenreList{
		Genres: genreList,
	}, nil
}
