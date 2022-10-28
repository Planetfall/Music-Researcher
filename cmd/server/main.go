package main

import (
	"fmt"
	"log"
	"net"

	flag "github.com/spf13/pflag"

	"github.com/Dadard29/planetfall/musicresearcher/internal/server"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	pb "github.com/Dadard29/planetfall/musicresearcher/pkg/pb"
)

func mustBindEnv(envName string) {
	viper.MustBindEnv(envName)
	if !viper.IsSet(envName) {
		log.Fatalf("Could not retrieve config env: %s\n", envName)
	}
}

func setConfig() {
	// from env
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("K_SERVICE", "music-researcher")

	mustBindEnv("PORT")
	mustBindEnv("K_SERVICE")
	mustBindEnv("SPOTIFY_CLIENT_ID")
	mustBindEnv("SPOTIFY_CLIENT_SECRET")

	// from cmd line
	flag.String("env", server.Production, "server environment")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)
}

func main() {
	setConfig()

	port := viper.GetString("PORT")

	musicResearcherPort := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp4", musicResearcherPort)
	if err != nil {
		log.Fatal(err)
	}
	musicResearcherServer := grpc.NewServer()

	log.Println("initializing server...")
	serv, err := server.NewServer(
		viper.GetString("env"),
		viper.GetString("K_SERVICE"),
		viper.GetString("SPOTIFY_CLIENT_ID"),
		viper.GetString("SPOTIFY_CLIENT_SECRET"),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer serv.Close()

	pb.RegisterMusicResearcherServer(musicResearcherServer, serv)

	log.Printf("listening on %s...\n", musicResearcherPort)
	if err := musicResearcherServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
