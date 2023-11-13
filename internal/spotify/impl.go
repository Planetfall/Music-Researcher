package myspotify

import (
	"context"
	"fmt"
	"log"
	"time"
)

type MySpotifyImpl struct {
	client          Client
	clientId        string
	clientSecret    string
	tokenExpiryTime time.Time

	provider Provider
	logger   *log.Logger
}

type MySpotifyOptions struct {
	ClientId     string
	ClientSecret string
	BaseLogger   *log.Logger
	Provider     Provider
}

func (opt MySpotifyOptions) getProvider() Provider {
	if opt.Provider == nil {
		return &providerImpl{}
	}

	return opt.Provider
}

func NewMySpotify(opt MySpotifyOptions) MySpotify {

	prefix := fmt.Sprintf("%s[%s] ", opt.BaseLogger.Prefix(), "SPOTIFY")
	logger := log.New(opt.BaseLogger.Writer(), prefix, opt.BaseLogger.Flags())

	return &MySpotifyImpl{
		client:          nil,
		clientId:        opt.ClientId,
		clientSecret:    opt.ClientSecret,
		tokenExpiryTime: time.Now(),

		provider: opt.getProvider(),
		logger:   logger,
	}
}

func (s *MySpotifyImpl) refresh(ctx context.Context) error {
	// current token is valid, no need to refresh
	if s.client != nil && s.tokenExpiryTime.After(time.Now()) {
		return nil
	}

	// current token expired, refreshing...
	s.logger.Println("refreshing spotify client...")

	client, expiryTime, err := s.provider.NewClient(ctx, s.clientId, s.clientSecret)
	if err != nil {
		return fmt.Errorf("spotify.refresh: %v", err)
	}

	s.client = client
	s.tokenExpiryTime = expiryTime

	// refreshed
	s.logger.Printf("spotify client refreshed, token expires in %v",
		time.Until(s.tokenExpiryTime))

	return nil
}
