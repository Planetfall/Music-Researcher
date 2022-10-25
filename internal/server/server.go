package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	pb "github.com/Dadard29/planetfall/musicresearcher/pkg/pb"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

// Main server struct
// holds clients configurations
type Server struct {
	pb.UnimplementedMusicResearcherServer
	metadataClient *metadata.Client
	projectID      string
	serviceName    string

	secretManager  *secretmanager.Client
	errorReporting *errorreporting.Client

	spotifyClient       *spotify.Client
	spotifyToken        *oauth2.Token
	spotifyClientID     string
	spotifyClientSecret string
}

func (s *Server) errorReport(err error, message string) {
	err = fmt.Errorf("%s: %v", message, err)
	s.errorReporting.Report(errorreporting.Entry{
		Error: err,
	})
	log.Println(err)
}

func (s *Server) Close() {
	s.secretManager.Close()
	s.errorReporting.Close()
}

func NewServer(
	serviceName string,
	spotifyClientID string,
	spotifyClientSecret string,
) (*Server, error) {

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
		metadataClient: metadataClient,
		projectID:      projectID,
		serviceName:    serviceName,

		secretManager:  secretManager,
		errorReporting: errorReporting,

		spotifyClient:       nil,
		spotifyToken:        nil,
		spotifyClientID:     spotifyClientID,
		spotifyClientSecret: spotifyClientSecret,
	}, nil
}
