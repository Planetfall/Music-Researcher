package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"

	pb "github.com/Dadard29/planetfall/music-researcher/musicresearcher"
)

type server struct {
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

func (s *server) errorReport(err error, message string) {
	err = fmt.Errorf("%s: %v", message, err)
	s.errorReporting.Report(errorreporting.Entry{
		Error: err,
	})
	log.Println(err)
}

func (s *server) close() {
	s.secretManager.Close()
	s.errorReporting.Close()
}

func newServer() (*server, error) {
	ctx := context.Background()

	var (
		serviceName         = os.Getenv("K_SERVICE")
		spotifyClientID     = os.Getenv("SPOTIFY_CLIENT_ID")
		spotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	)

	// init metadata client
	metadataClient := metadata.NewClient(&http.Client{})
	projectID, err := metadataClient.ProjectID()
	if err != nil {
		return nil, err
	}

	// init secret manager
	secretManager, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	// init error reporting
	errorReporting, err := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create error reporting: %v", err)
	}

	return &server{
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	musicResearcherPort := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp4", musicResearcherPort)
	if err != nil {
		log.Fatal(err)
	}
	musicResearcherServer := grpc.NewServer()

	serv, err := newServer()
	if err != nil {
		log.Fatal(err)
	}

	defer serv.close()

	pb.RegisterMusicResearcherServer(musicResearcherServer, serv)

	if err := musicResearcherServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
