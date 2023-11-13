package runserver

import (
	"fmt"
	"log"
	"net"

	"github.com/planetfall/framework/pkg/config"
	"github.com/planetfall/framework/pkg/server"
	"github.com/planetfall/musicresearcher/internal/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	portFlag                = "port"
	serviceFlag             = "service"
	spotifyClientIdFlag     = "spotify-client-id"
	spotifyClientSecretFlag = "spotify-client-secret"
)

func getConfig() (config.Config, error) {

	entries := []config.Entry{{
		Flag:         portFlag,
		DefaultValue: "8080",
		Description:  "the exposed port of the service",
		EnvKey:       "PORT",
	}, {
		Flag:         serviceFlag,
		DefaultValue: "cloud-microservice",
		Description:  "the service name",
		EnvKey:       "K_SERVICE",
	}, {
		Flag:         spotifyClientIdFlag,
		DefaultValue: "",
		Description:  "the client ID for Spotify OAuth authentication",
		EnvKey:       "SPOTIFY_CLIENT_ID",
	}, {
		Flag:         spotifyClientSecretFlag,
		DefaultValue: "",
		Description:  "the client secret for Spotify OAuth authentication",
		EnvKey:       "SPOTIFY_CLIENT_SECRET",
	},
	}

	c, err := config.NewConfig(entries)
	if err != nil {
		return nil, fmt.Errorf("config.NewConfig: %v", err)
	}

	return c, nil
}

func getServer(cfg config.Config) (*server.Server, error) {

	serviceName := viper.GetString(serviceFlag)
	s, err := server.NewServer(cfg, serviceName)
	if err != nil {
		return nil, fmt.Errorf("server.NewServer: %v", err)
	}

	return s, nil
}

func getListener() (net.Listener, error) {
	port := viper.GetString(portFlag)
	addr := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, fmt.Errorf("net.Listen: %v", err)
	}

	return lis, nil
}

func getSpotifyCredentials() (string, string, error) {
	spotifyClientId := viper.GetString(spotifyClientIdFlag)
	if spotifyClientId == "" {
		return "", "", fmt.Errorf("spotify client ID not provided")
	}

	spotifyClientSecret := viper.GetString(spotifyClientSecretFlag)
	if spotifyClientSecret == "" {
		return "", "", fmt.Errorf("spotify client secret not provided")
	}

	return spotifyClientId, spotifyClientSecret, nil
}

func RunServer() {
	log.SetPrefix("[RUNSERVER] ")

	// config
	log.Printf("setting up the config")
	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("getConfig: %v", err)
	}

	// server
	log.Printf("setting up the server")
	srv, err := getServer(cfg)
	if err != nil {
		log.Fatalf("getServer: %v", err)
	}

	// service
	log.Printf("setting up the service")
	grpc := grpc.NewServer()
	spotifyClientId, spotifyClientSecret, err := getSpotifyCredentials()
	if err != nil {
		log.Fatalf("getSpotifyCredentials: %v", err)
	}
	svc := service.NewService(
		grpc, srv, spotifyClientId, spotifyClientSecret)

	// service start
	lis, err := getListener()
	if err != nil {
		log.Fatalf("getListener: %v", err)
	}

	log.Printf("starting listening on %s", lis.Addr().String())
	if err := svc.Start(lis); err != nil {
		log.Fatalf("svc.Start: %v", err)
	}

	log.Printf("service stopped")
}
