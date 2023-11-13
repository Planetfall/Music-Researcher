package service

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/planetfall/framework/pkg/server"
	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
	myspotify "github.com/planetfall/musicresearcher/internal/spotify"
	"google.golang.org/grpc"
)

type Service struct {
	pb.UnimplementedMusicResearcherServer
	grpcSrv *grpc.Server

	srv *server.Server

	mySpotify myspotify.MySpotify
}

func NewService(
	grpcSrv *grpc.Server,
	srv *server.Server,

	spotifyClientId string,
	spotifyClientSecret string,
) *Service {

	mySpotify := myspotify.NewMySpotify(myspotify.MySpotifyOptions{
		ClientId:     spotifyClientId,
		ClientSecret: spotifyClientSecret,
		BaseLogger:   srv.Logger,
	})

	newService := &Service{
		grpcSrv: grpcSrv,
		srv:     srv,

		mySpotify: mySpotify,
	}

	pb.RegisterMusicResearcherServer(grpcSrv, newService)

	return newService
}

func (s *Service) Start(lis net.Listener) error {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.grpcSrv.Serve(lis); err != nil {
			log.Fatalf("grpc.Serve: %v", err)
		}
	}()

	<-done

	defer s.srv.Close()
	s.grpcSrv.GracefulStop()

	return nil
}
