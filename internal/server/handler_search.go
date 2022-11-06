package server

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Dadard29/planetfall/musicresearcher/pkg/pb"
	"github.com/zmb3/spotify/v2"
)

func (s *Server) listArtistsFromTrack(ctx context.Context, track spotify.FullTrack,
	artistBufferList []spotify.FullArtist) ([]spotify.FullArtist, error) {

	out := make([]spotify.FullArtist, 0)
	for _, artist := range track.Artists {
		// check if track artist already in buffer
		inBuffer := false
		for _, artistBuffer := range artistBufferList {
			if artist.ID == artistBuffer.ID {
				// store it in buffer, to avoid requesting it again
				out = append(out, artistBuffer)
				inBuffer = true
			}
		}

		// if not, request the artist from API
		if !inBuffer {
			artistBuffer, err := s.spotifyClient.GetArtist(ctx, artist.ID)
			if err != nil {
				return nil, err
			}
			out = append(out, *artistBuffer)
		}
	}

	return out, nil
}

func (s *Server) pagesToTrackList(
	ctx context.Context, pages *spotify.FullTrackPage) ([]*pb.Track, error) {

	var trackList = make([]*pb.Track, 0)
	var artistBufferList = make([]spotify.FullArtist, 0)

	for {
		for _, track := range pages.Tracks {
			artistList, err := s.listArtistsFromTrack(ctx, track, artistBufferList)
			if err != nil {
				return nil, err
			}

			track := s.newTrack(track, artistList)
			trackList = append(trackList, track)
		}

		if err := s.spotifyClient.NextPage(ctx, pages); err == spotify.ErrNoMorePages {
			break
		}
	}

	return trackList, nil
}

func (s *Server) Search(ctx context.Context, params *pb.Parameters) (*pb.Results, error) {

	err := s.setSpotifyClient(ctx)
	if err != nil {
		s.errorReport(err, "failed to setup connection with Spotify")
		return nil, err
	}

	// query
	query := params.Query
	if query == "" {
		return nil, fmt.Errorf("query parameter is empty")
	}

	// genre list
	// if genre list supplied, add genre filters to the query
	if len(params.GenreFilters) > 0 {
		genreFilters := ""
		for _, genre := range params.GenreFilters {
			genreFilters = fmt.Sprintf("%s genre:%s", genreFilters, genre)
		}
		query = fmt.Sprintf("%s %s", query, genreFilters)
	}

	// limit
	const defaultLimit = 10
	limit := int(params.Limit)
	if limit <= 0 {
		limit = defaultLimit
	}

	log.Printf("requesting Spotify with query: %s", query)
	searchResult, err := s.spotifyClient.Search(
		ctx, query, spotify.SearchTypeTrack, spotify.Limit(limit))
	if err != nil {
		s.errorReport(err, "failed interacting with Spotify")
		return nil, err
	}

	pages := searchResult.Tracks
	trackList, err := s.pagesToTrackList(ctx, pages)
	if err != nil {
		s.errorReport(err, "failed to extract track pages")
		return nil, err
	}

	return &pb.Results{
		Albums:  nil,
		Artists: nil,
		Tracks:  trackList,
	}, nil
}
