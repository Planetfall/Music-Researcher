package myspotify

import (
	"context"
	"net/http"

	"github.com/zmb3/spotify/v2"
)

type Client interface {
	GetAvailableGenreSeeds(ctx context.Context) ([]string, error)

	GetArtist(ctx context.Context,
		artistID spotify.ID) (*spotify.FullArtist, error)

	NextPage(ctx context.Context, p *spotify.FullTrackPage) error

	Search(ctx context.Context, query string,
		t spotify.SearchType,
		opts ...spotify.RequestOption) (*spotify.SearchResult, error)
}

// clientImpl is a wrapper around spotify.Client
// It is needed as the method NextPage from the interface
// is implemented using an internal type (pageable).
// Thus, NextPage from the clientImpl is wrapped using the exported
// spotify.FullTrackPage.
type clientImpl struct {
	spotifyClient *spotify.Client
}

func newClientImpl(h *http.Client) *clientImpl {
	return &clientImpl{
		spotifyClient: spotify.New(h),
	}
}

func (c *clientImpl) GetArtist(ctx context.Context,
	artistID spotify.ID) (*spotify.FullArtist, error) {

	return c.spotifyClient.GetArtist(ctx, artistID)
}

func (c *clientImpl) GetAvailableGenreSeeds(
	ctx context.Context) ([]string, error) {

	return c.spotifyClient.GetAvailableGenreSeeds(ctx)
}

func (c *clientImpl) NextPage(
	ctx context.Context, p *spotify.FullTrackPage) error {

	return c.spotifyClient.NextPage(ctx, p)
}

func (c *clientImpl) Search(ctx context.Context, query string,
	t spotify.SearchType,
	opts ...spotify.RequestOption) (*spotify.SearchResult, error) {

	return c.spotifyClient.Search(ctx, query, t, opts...)
}
