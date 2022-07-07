package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"time"

	pb "github.com/Dadard29/planetfall/music-researcher/musicresearcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var query = flag.String("query", "chilly gonzales", "The query to send")
var genre = flag.String("genre", "", "The genre to filter")

func main() {
	flag.Parse()
	log.Println(*query, *genre)

	host := "music-researcher-twecq3u42q-ew.a.run.app:443"

	var opts []grpc.DialOption
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Println("failed getting certs")
		log.Fatal(err)
	}

	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))

	conn, err := grpc.Dial(host, opts...)
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
		Limit:        10,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(r)

}
