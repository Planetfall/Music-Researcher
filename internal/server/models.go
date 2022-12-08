package server

import (
	pb "github.com/Dadard29/planetfall/musicresearcher/pkg/pb"
	"github.com/zmb3/spotify/v2"
)

const spotifyUrlKey = "spotify"

type ItemType int

// fixme: https://github.com/Dadard29/planetfall/issues/4
var defaultImageUrl = map[pb.Type]string{
	pb.Type_UNKNOWN: "",
	pb.Type_ARTIST:  "",
	pb.Type_ALBUM:   "",
	pb.Type_TRACK:   "",
}

func (s *Server) getImageUrl(images []spotify.Image, itemType pb.Type) string {
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

func (s *Server) newTrack(track spotify.FullTrack, artistList []spotify.FullArtist) *pb.Track {
	albumDto := &pb.Album{
		ID:          track.Album.ID.String(),
		Name:        track.Album.Name,
		ReleaseDate: track.Album.ReleaseDate,
		SpotifyUrl:  track.Album.ExternalURLs[spotifyUrlKey],
		ImageUrl:    s.getImageUrl(track.Album.Images, pb.Type_ALBUM),
	}

	artistDtoList := make([]*pb.Artist, 0)
	for _, artist := range artistList {
		artistDtoList = append(artistDtoList, &pb.Artist{
			ID:         artist.ID.String(),
			Name:       artist.Name,
			SpotifyUrl: artist.ExternalURLs[spotifyUrlKey],
			Genres:     artist.Genres,
			ImageUrl:   s.getImageUrl(artist.Images, pb.Type_ARTIST),
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
