package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
)

const (
	Development string = "development"
	Production  string = "production"
)

func newServerDevelopement(
	serviceName string,
	spotifyClientID string,
	spotifyClientSecret string,
) (*Server, error) {
	log.Printf("setting up server for %s...\n", Development)

	log.SetPrefix("[DEV] ")

	return &Server{
		env:                 Development,
		serviceName:         serviceName,
		metadataClient:      nil,
		projectID:           "",
		secretManager:       nil,
		errorReporting:      nil,
		spotifyClient:       nil,
		spotifyToken:        nil,
		spotifyClientID:     spotifyClientID,
		spotifyClientSecret: spotifyClientSecret,
	}, nil
}

func newServerProduction(
	serviceName string,
	spotifyClientID string,
	spotifyClientSecret string,
) (*Server, error) {
	log.Printf("setting up server for %s...\n", Production)

	log.SetPrefix("[PRD] ")

	ctx := context.Background()

	// init metadata client
	log.Println("initializing metadata client...")
	metadataClient := metadata.NewClient(&http.Client{})
	projectID, err := metadataClient.ProjectID()
	if err != nil {
		return nil, err
	}

	// init secret manager
	log.Println("initializing secret manager...")
	secretManager, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	// init error reporting
	log.Println("initializing error reporting...")
	errorReporting, err := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create error reporting: %v", err)
	}

	return &Server{
		env:         Production,
		serviceName: serviceName,

		metadataClient: metadataClient,
		projectID:      projectID,
		secretManager:  secretManager,
		errorReporting: errorReporting,

		spotifyClient:       nil,
		spotifyToken:        nil,
		spotifyClientID:     spotifyClientID,
		spotifyClientSecret: spotifyClientSecret,
	}, nil
}
