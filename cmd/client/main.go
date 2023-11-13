package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"time"

	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var query = flag.String("query", "chilly gonzales", "The query to send")
var genre = flag.String("genre", "", "The genre to filter")
var host = flag.String("host", "music-researcher-twecq3u42q-ew.a.run.app:443", "The service's host")
var use_tls = flag.Bool("tls", true, "use tls")

func main() {
	flag.Parse()
	log.Println(*host, *query, *genre)

	var opts []grpc.DialOption
	if *use_tls {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			log.Println("failed getting certs")
			log.Fatal(err)
		}

		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})

		opts = append(opts, grpc.WithTransportCredentials(cred))
	} else {
		opts = append(opts,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(*host, opts...)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c := pb.NewMusicResearcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := c.Search(ctx, &pb.Parameters{
		Query:        *query,
		GenreFilters: []string{*genre},
		Limit:        3,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(r)
	log.Println(len(r.Tracks))

}
