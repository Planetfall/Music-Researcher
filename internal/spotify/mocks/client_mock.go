package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/zmb3/spotify/v2"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) GetAvailableGenreSeeds(
	ctx context.Context) ([]string, error) {

	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *ClientMock) GetArtist(ctx context.Context,
	artistID spotify.ID) (*spotify.FullArtist, error) {

	args := m.Called(artistID)
	return args.Get(0).(*spotify.FullArtist), args.Error(1)
}

func (m *ClientMock) NextPage(
	ctx context.Context, p *spotify.FullTrackPage) error {

	args := m.Called()
	return args.Error(0)
}

func (m *ClientMock) Search(ctx context.Context, query string,
	t spotify.SearchType,
	opts ...spotify.RequestOption) (*spotify.SearchResult, error) {

	args := m.Called(query)
	return args.Get(0).(*spotify.SearchResult), args.Error(1)
}
