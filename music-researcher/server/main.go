package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/errorreporting"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/xds"

	pb "github.com/Dadard29/planetfall/music-researcher/musicresearcher"
)

const (
	spotifyClientID     = "SPOTIFY_CLIENT_ID"
	spotifyClientSecret = "SPOTIFY_CLIENT_SECRET"
)

var projectID = os.Getenv("PROJECT_ID")
var serviceName = os.Getenv("SERVICE")

type server struct {
	pb.UnimplementedMusicResearcherServer

	secretManager  *secretmanager.Client
	errorReporting *errorreporting.Client

	spotifyClient *spotify.Client
	spotifyToken  *oauth2.Token
}

func (s *server) getSecret(secretName string) (string, error) {
	ctx := context.Background()

	secretPath := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}

	result, err := s.secretManager.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}
	return string(result.Payload.Data), nil
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
		secretManager:  secretManager,
		errorReporting: errorReporting,

		spotifyClient: nil,
		spotifyToken:  nil,
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
	creds := insecure.NewCredentials()
	musicResearcherServer := xds.NewGRPCServer(grpc.Creds(creds))

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
