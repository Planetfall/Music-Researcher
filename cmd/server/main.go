package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Dadard29/planetfall/musicresearcher/internal/server"
	"google.golang.org/grpc"

	pb "github.com/Dadard29/planetfall/musicresearcher/pkg/pb"
)

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

	log.Println("initializing server...")
	serv, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	defer serv.Close()

	pb.RegisterMusicResearcherServer(musicResearcherServer, serv)

	log.Printf("listening on %s\n...", musicResearcherPort)
	if err := musicResearcherServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
