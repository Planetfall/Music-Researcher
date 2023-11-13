package myspotify

import (
	"context"
	"fmt"

	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
	"github.com/zmb3/spotify/v2"
)

const (
	defaultSearchLimit = 10
)

func addGenreListToQuery(query string, genreFilters []string) string {
	if len(genreFilters) > 0 {
		for _, genre := range genreFilters {
			query = fmt.Sprintf("%s genre:%s", query, genre)
		}
	}

	return query
}

func (s *MySpotifyImpl) Search(ctx context.Context,
	query string, genreFilters []string, limit int) ([]*pb.Track, error) {

	if err := s.refresh(ctx); err != nil {
		return nil, err
	}

	// validate limit
	if limit <= 0 {
		limit = defaultSearchLimit
	}

	// validate query
	if query == "" {
		return nil, fmt.Errorf("provided query is empty")
	}

	// format query with genre list
	query = addGenreListToQuery(query, genreFilters)

	// performs the search
	s.logger.Printf("querying spotify with query `%v`", query)
	results, err := s.client.Search(ctx, query, spotify.SearchTypeTrack, spotify.Limit(limit))
	if err != nil {
		return nil, fmt.Errorf("client.Search: %v", err)
	}

	trackList, err := s.pagesToTrackList(ctx, results.Tracks)
	if err != nil {
		return nil, fmt.Errorf("pagesToTrackList: %v", err)
	}

	return trackList, nil
}
