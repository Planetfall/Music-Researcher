package myspotify

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/endpoints"

	spotifyAuth "github.com/zmb3/spotify/v2/auth"
)

type Provider interface {
	NewClient(ctx context.Context,
		clientId string,
		clientSecret string) (Client, time.Time, error)
}

type providerImpl struct {
}

func (p *providerImpl) NewClient(ctx context.Context,
	clientId string, clientSecret string) (Client, time.Time, error) {

	oauthConfig := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     endpoints.Spotify.TokenURL,
	}

	token, err := oauthConfig.Token(ctx)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to get spotify oauth token: %v", err)
	}

	httpClient := spotifyAuth.New().Client(ctx, token)
	client := newClientImpl(httpClient)

	return client, token.Expiry, nil
}
