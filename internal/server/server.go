package server

import (
	"fmt"
	"log"

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
	env         string
	serviceName string

	metadataClient *metadata.Client
	projectID      string
	secretManager  *secretmanager.Client
	errorReporting *errorreporting.Client

	spotifyClient       *spotify.Client
	spotifyToken        *oauth2.Token
	spotifyClientID     string
	spotifyClientSecret string
}

func (s *Server) errorReport(err error, message string) {
	err = fmt.Errorf("%s: %v", message, err)
	if s.errorReporting != nil {
		s.errorReporting.Report(errorreporting.Entry{
			Error: err,
		})
	}
	log.Println(err)
}

func (s *Server) Close() {
	if s.secretManager != nil {
		s.secretManager.Close()
	}

	if s.errorReporting != nil {
		s.errorReporting.Close()
	}
}

func NewServer(
	env string,
	serviceName string,
	spotifyClientID string,
	spotifyClientSecret string,
) (*Server, error) {

	switch env {
	case Development:
		return newServerDevelopement(
			serviceName, spotifyClientID, spotifyClientSecret)
	case Production:
		return newServerProduction(
			serviceName, spotifyClientID, spotifyClientSecret)
	default:
		return nil, fmt.Errorf(
			"failed to create server with unsupported environment: %s\n", env)
	}
}
