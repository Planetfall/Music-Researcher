package myspotify

import (
	"context"

	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
	"github.com/zmb3/spotify/v2"
)

const spotifyUrlKey = "spotify"

type ItemType int

// fixme: https://github.com/planetfall/issues/4
var defaultImageUrl = map[pb.Type]string{
	pb.Type_UNKNOWN: "",
	pb.Type_ARTIST:  "",
	pb.Type_ALBUM:   "",
	pb.Type_TRACK:   "",
}

func getImageUrl(images []spotify.Image, itemType pb.Type) string {
	// if no image found in spotify metadatas
	if len(images) == 0 {
		if url, check := defaultImageUrl[itemType]; !check {
			return defaultImageUrl[pb.Type_UNKNOWN]
		} else {
			return url
		}
	}

	return images[0].URL
}

func mapSpotifyTrack(track spotify.FullTrack, artistList []spotify.FullArtist) *pb.Track {
	albumDto := &pb.Album{
		ID:          track.Album.ID.String(),
		Name:        track.Album.Name,
		ReleaseDate: track.Album.ReleaseDate,
		SpotifyUrl:  track.Album.ExternalURLs[spotifyUrlKey],
		ImageUrl:    getImageUrl(track.Album.Images, pb.Type_ALBUM),
	}

	artistDtoList := make([]*pb.Artist, 0)
	for _, artist := range artistList {
		artistDtoList = append(artistDtoList, &pb.Artist{
			ID:         artist.ID.String(),
			Name:       artist.Name,
			SpotifyUrl: artist.ExternalURLs[spotifyUrlKey],
			Genres:     artist.Genres,
			ImageUrl:   getImageUrl(artist.Images, pb.Type_ARTIST),
		})
	}

	trackDto := &pb.Track{
		ID:         track.ID.String(),
		Name:       track.Name,
		SpotifyUrl: track.ExternalURLs[spotifyUrlKey],
		DurationMs: int32(track.Duration),
		PreviewUrl: track.PreviewURL,
		Popularity: int32(track.Popularity),

		Album:   albumDto,
		Artists: artistDtoList,
	}

	return trackDto
}

// lists the full artists metadatas from a full spotify trac
func (s *MySpotifyImpl) listArtistsFromTrack(ctx context.Context,
	track spotify.FullTrack, artistBufferList []spotify.FullArtist,
) ([]spotify.FullArtist, error) {

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
			artistBuffer, err := s.client.GetArtist(ctx, artist.ID)
			if err != nil {
				return nil, err
			}
			out = append(out, *artistBuffer)
		}
	}

	return out, nil
}

// converts a list of track pages into the output format
// and enrich the result with the full artist metadatas
func (s *MySpotifyImpl) pagesToTrackList(
	ctx context.Context, pages *spotify.FullTrackPage) ([]*pb.Track, error) {

	var trackList = make([]*pb.Track, 0)
	var artistBufferList = make([]spotify.FullArtist, 0)

	for {
		for _, track := range pages.Tracks {
			artistList, err := s.listArtistsFromTrack(ctx, track, artistBufferList)
			if err != nil {
				return nil, err
			}

			track := mapSpotifyTrack(track, artistList)
			trackList = append(trackList, track)
		}

		if err := s.client.NextPage(ctx, pages); err == spotify.ErrNoMorePages {
			break
		}
	}

	return trackList, nil
}
