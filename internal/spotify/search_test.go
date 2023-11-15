package myspotify_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	myspotify "github.com/planetfall/musicresearcher/internal/spotify"
	"github.com/planetfall/musicresearcher/internal/spotify/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
)

func getArtist(artistId spotify.ID) *spotify.FullArtist {
	return &spotify.FullArtist{
		SimpleArtist: spotify.SimpleArtist{
			ID:   artistId,
			Name: "artist name",
			ExternalURLs: map[string]string{
				"spotify": "artist-spotify-url",
			},
		},
		Genres: []string{"genre1"},
	}
}

func getSearchResults(artistId spotify.ID) *spotify.SearchResult {

	return &spotify.SearchResult{
		Tracks: &spotify.FullTrackPage{
			Tracks: []spotify.FullTrack{
				{
					Album: spotify.SimpleAlbum{
						ID:          "album 1",
						Name:        "album",
						ReleaseDate: "2023-10-01",
						ExternalURLs: map[string]string{
							"spotify": "album-spotify-url",
						},
					},
					SimpleTrack: spotify.SimpleTrack{
						Artists: []spotify.SimpleArtist{
							{ID: artistId},
						},
					},
				},
				{
					Album: spotify.SimpleAlbum{
						ID:          "album 2",
						Name:        "album",
						ReleaseDate: "2023-10-01",
						ExternalURLs: map[string]string{
							"spotify": "album-spotify-url",
						},
					},
					SimpleTrack: spotify.SimpleTrack{
						Artists: []spotify.SimpleArtist{
							{ID: artistId},
						},
					},
				},
			},
		},
	}

}

func newMySpotifyClient(client myspotify.Client) myspotify.MySpotify {

	clientIdGiven := "client-id"
	clientSecretGiven := "client-secret"
	expiryGiven := time.Now().Add(time.Second * 5)

	providerGiven := &mocks.ProviderMock{}
	providerGiven.
		On("NewClient", clientIdGiven, clientSecretGiven).
		Return(client, expiryGiven, nil)

	optGiven := myspotify.MySpotifyOptions{
		ClientId:     clientIdGiven,
		ClientSecret: clientSecretGiven,
		BaseLogger:   log.Default(),
		Provider:     providerGiven,
	}

	return myspotify.NewMySpotify(optGiven)
}

func TestSearch(t *testing.T) {

	ctxGiven := context.Background()
	queryGiven := "chilly gonzales crying"
	genreListGiven := []string{}
	limitGiven := 10

	artistIdGiven := spotify.ID("artist-id-1")
	searchResultsGiven := getSearchResults(artistIdGiven)
	artistGiven := getArtist(artistIdGiven)

	clientGiven := &mocks.ClientMock{}
	clientGiven.On("Search", queryGiven).Return(searchResultsGiven, nil)
	clientGiven.On("GetArtist", artistIdGiven).Return(artistGiven, nil)
	clientGiven.On("NextPage").Return(spotify.ErrNoMorePages)

	mySpotifyClient := newMySpotifyClient(clientGiven)
	results, err := mySpotifyClient.Search(
		ctxGiven, queryGiven, genreListGiven, limitGiven)
	assert.Nil(t, err)
	assert.Len(t, results, 2)

	clientGiven.AssertExpectations(t)
}

func TestSearch_withGenreFilters(t *testing.T) {

	ctxGiven := context.Background()
	queryGiven := "chilly gonzales crying"
	genreListGiven := []string{"genre1"}
	queryWithGenresGiven := "chilly gonzales crying genre:genre1"
	limitGiven := 10

	artistIdGiven := spotify.ID("artist-id-1")
	searchResultsGiven := getSearchResults(artistIdGiven)
	artistGiven := getArtist(artistIdGiven)

	clientGiven := &mocks.ClientMock{}
	clientGiven.On("Search", queryWithGenresGiven).Return(searchResultsGiven, nil)
	clientGiven.On("GetArtist", artistIdGiven).Return(artistGiven, nil)
	clientGiven.On("NextPage").Return(spotify.ErrNoMorePages)

	mySpotifyClient := newMySpotifyClient(clientGiven)
	results, err := mySpotifyClient.Search(
		ctxGiven, queryGiven, genreListGiven, limitGiven)
	assert.Nil(t, err)
	assert.Len(t, results, 2)

	clientGiven.AssertExpectations(t)
}

func TestSearch_withGetArtistError(t *testing.T) {

	ctxGiven := context.Background()
	queryGiven := "chilly gonzales crying"
	genreListGiven := []string{}
	limitGiven := 10

	artistIdGiven := spotify.ID("artist-id-1")
	searchResultsGiven := getSearchResults(artistIdGiven)
	errorGiven := fmt.Errorf("failed to get artist")

	clientGiven := &mocks.ClientMock{}
	clientGiven.On("Search", queryGiven).Return(searchResultsGiven, nil)
	clientGiven.
		On("GetArtist", artistIdGiven).
		Return(&spotify.FullArtist{}, errorGiven)

	mySpotifyClient := newMySpotifyClient(clientGiven)
	results, err := mySpotifyClient.Search(
		ctxGiven, queryGiven, genreListGiven, limitGiven)
	assert.Nil(t, results)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "client.GetArtist")

	clientGiven.AssertExpectations(t)
}

func TestSearch_withEmptyQuery(t *testing.T) {

	ctxGiven := context.Background()
	queryGiven := ""
	genreListGiven := []string{}
	limitGiven := 10

	clientGiven := &mocks.ClientMock{}
	mySpotifyClient := newMySpotifyClient(clientGiven)
	results, err := mySpotifyClient.Search(
		ctxGiven, queryGiven, genreListGiven, limitGiven)
	assert.NotNil(t, err)
	assert.Nil(t, results)
}

func TestSearch_withNewClientError(t *testing.T) {

	ctxGiven := context.Background()
	queryGiven := ""
	genreListGiven := []string{}
	limitGiven := 10

	clientIdGiven := "client-id"
	clientSecretGiven := "client-secret"
	errorGiven := fmt.Errorf("failed to create client")

	providerGiven := &mocks.ProviderMock{}
	providerGiven.
		On("NewClient", clientIdGiven, clientSecretGiven).
		Return(&mocks.ClientMock{}, time.Now(), errorGiven)

	optGiven := myspotify.MySpotifyOptions{
		ClientId:     clientIdGiven,
		ClientSecret: clientSecretGiven,
		BaseLogger:   log.Default(),
		Provider:     providerGiven,
	}

	mySpotifyClient := myspotify.NewMySpotify(optGiven)
	results, err := mySpotifyClient.Search(
		ctxGiven, queryGiven, genreListGiven, limitGiven)
	assert.Nil(t, results)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "provider.NewClient")

	providerGiven.AssertExpectations(t)
}
