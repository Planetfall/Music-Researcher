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
)

func TestGetGenreList(t *testing.T) {

	genreListGiven := []string{"genre1", "genre2"}

	ctxGiven := context.Background()

	clientGiven := &mocks.ClientMock{}
	clientGiven.On("GetAvailableGenreSeeds").Return(genreListGiven, nil)

	mySpotifyClient := newMySpotifyClient(clientGiven)

	genreListResponse, err := mySpotifyClient.GetGenreList(ctxGiven)
	assert.Nil(t, err)

	genreListActual := genreListResponse.Genres
	assert.Equal(t, genreListActual, genreListGiven)

	clientGiven.AssertExpectations(t)
}

func TestGetGenreList_withoutRefresh(t *testing.T) {

	clientIdGiven := "client-id"
	clientSecretGiven := "client-secret"
	expiryGiven := time.Now().Add(time.Second * 5)
	genreListGiven := []string{"genre1", "genre2"}

	ctxGiven := context.Background()

	clientGiven := &mocks.ClientMock{}
	clientGiven.On("GetAvailableGenreSeeds").Return(genreListGiven, nil)

	providerGiven := &mocks.ProviderMock{}
	providerGiven.
		On("NewClient", clientIdGiven, clientSecretGiven).
		Return(clientGiven, expiryGiven, nil)

	optGiven := myspotify.MySpotifyOptions{
		ClientId:     clientIdGiven,
		ClientSecret: clientSecretGiven,
		BaseLogger:   log.Default(),
		Provider:     providerGiven,
	}

	mySpotifyClient := myspotify.NewMySpotify(optGiven)

	// first call
	genreListResponse, err := mySpotifyClient.GetGenreList(ctxGiven)
	assert.Nil(t, err)

	genreListActual := genreListResponse.Genres
	assert.Equal(t, genreListActual, genreListGiven)

	time.Sleep(2 * time.Second)

	// second call
	genreListResponse, err = mySpotifyClient.GetGenreList(ctxGiven)
	assert.Nil(t, err)

	genreListActual = genreListResponse.Genres
	assert.Equal(t, genreListActual, genreListGiven)

	providerGiven.AssertExpectations(t)
	providerGiven.AssertNumberOfCalls(t, "NewClient", 1)

	clientGiven.AssertExpectations(t)
	clientGiven.AssertNumberOfCalls(t, "GetAvailableGenreSeeds", 2)
}

func TestGetGenreList_withRefresh(t *testing.T) {

	clientIdGiven := "client-id"
	clientSecretGiven := "client-secret"
	expiryGiven := time.Now().Add(time.Second * 5)
	genreListGiven := []string{"genre1", "genre2"}

	ctxGiven := context.Background()

	clientGiven := &mocks.ClientMock{}
	clientGiven.On("GetAvailableGenreSeeds").Return(genreListGiven, nil)

	providerGiven := &mocks.ProviderMock{}
	providerGiven.
		On("NewClient", clientIdGiven, clientSecretGiven).
		Return(clientGiven, expiryGiven, nil)

	optGiven := myspotify.MySpotifyOptions{
		ClientId:     clientIdGiven,
		ClientSecret: clientSecretGiven,
		BaseLogger:   log.Default(),
		Provider:     providerGiven,
	}

	mySpotifyClient := myspotify.NewMySpotify(optGiven)

	// first call
	genreListResponse, err := mySpotifyClient.GetGenreList(ctxGiven)
	assert.Nil(t, err)

	genreListActual := genreListResponse.Genres
	assert.Equal(t, genreListActual, genreListGiven)

	time.Sleep(6 * time.Second)

	// second call
	genreListResponse, err = mySpotifyClient.GetGenreList(ctxGiven)
	assert.Nil(t, err)

	genreListActual = genreListResponse.Genres
	assert.Equal(t, genreListActual, genreListGiven)

	providerGiven.AssertExpectations(t)
	providerGiven.AssertNumberOfCalls(t, "NewClient", 2)

	clientGiven.AssertExpectations(t)
	clientGiven.AssertNumberOfCalls(t, "GetAvailableGenreSeeds", 2)
}

func TestGetGenreList_withProviderError(t *testing.T) {

	clientIdGiven := "client-id"
	clientSecretGiven := "client-secret"

	ctxGiven := context.Background()
	errorGiven := fmt.Errorf("failed to provide client")

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
	_, err := mySpotifyClient.GetGenreList(ctxGiven)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "provider.NewClient")

	providerGiven.AssertExpectations(t)
}

func TestGetGenreList_withClientError(t *testing.T) {
	clientIdGiven := "client-id"
	clientSecretGiven := "client-secret"
	expiryGiven := time.Now().Add(time.Second * 5)

	ctxGiven := context.Background()
	errorGiven := fmt.Errorf("failed to get available genre seeds")

	clientGiven := &mocks.ClientMock{}
	clientGiven.On("GetAvailableGenreSeeds").Return([]string{}, errorGiven)

	providerGiven := &mocks.ProviderMock{}
	providerGiven.
		On("NewClient", clientIdGiven, clientSecretGiven).
		Return(clientGiven, expiryGiven, nil)

	optGiven := myspotify.MySpotifyOptions{
		ClientId:     clientIdGiven,
		ClientSecret: clientSecretGiven,
		BaseLogger:   log.Default(),
		Provider:     providerGiven,
	}

	mySpotifyClient := myspotify.NewMySpotify(optGiven)
	_, err := mySpotifyClient.GetGenreList(ctxGiven)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "client.GetAvailableGenreSeeds")

	providerGiven.AssertExpectations(t)
	clientGiven.AssertExpectations(t)
}
