package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/endpoints"

	"github.com/zmb3/spotify/v2"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
)

func (s *server) setSpotifyClient(ctx context.Context) error {
	// current token is valid, no need to refresh
	if s.spotifyToken != nil && s.spotifyToken.Expiry.After(time.Now()) {
		return nil
	}

	// current token expired, refreshing...
	log.Println("refreshing spotify client...")

	clientId, err := s.getSecret(spotifyClientID)
	if err != nil {
		return fmt.Errorf("failed to get client ID: %v", err)
	}

	clientSecret, err := s.getSecret(spotifyClientSecret)
	if err != nil {
		return fmt.Errorf("failed to get client secret: %v", err)
	}

	oauthConfig := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     endpoints.Spotify.TokenURL,
	}

	token, err := oauthConfig.Token(ctx)
	if err != nil {
		return fmt.Errorf("failed to get spotify oauth token: %v", err)
	}

	s.spotifyToken = token

	httpClient := spotifyAuth.New().Client(ctx, token)
	s.spotifyClient = spotify.New(httpClient)

	return nil
}
