package myspotify

import (
	"context"

	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
)

type MySpotify interface {
	Search(ctx context.Context,
		query string, genreFilters []string,
		limit int) ([]*pb.Track, error)

	GetGenreList(ctx context.Context) (*pb.GenreList, error)
}
