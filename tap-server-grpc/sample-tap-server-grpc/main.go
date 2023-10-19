package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"

	"google.golang.org/grpc"
)

var (
	GrpcPort     = flag.Int("p", 8080, "port")
	dataScrubber data_scrubber.DataScrubber
)

type server struct{}

func (s *server) ReportTap(srv tap_service.TapService_ReportTapServer) error {
	log.Printf("Starting to listen for requests")
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			log.Printf("End of stream\n")
			return ctx.Err()
		default:
		}

		tapRequest, err := srv.Recv()
		if err == io.EOF {
			// Client has closed the stream
			return nil
		}
		log.Printf("got a request!")
		dataScrubber.ScrubTapRequest(tapRequest)
		tapRequestJson, err := json.MarshalIndent(tapRequest, "", "  ")
		if err != nil {
			log.Printf("Error marshalling proto message to json: %s", err.Error())
		}
		log.Printf("Message contents were: %s\n", tapRequestJson)
	}
}

func main() {
	flag.Parse()
	dataScrubber.Init()

	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(1000)}
	s := grpc.NewServer(sopts...)
	tap_service.RegisterTapServiceServer(s, &server{})

	address := fmt.Sprintf(":%d", *GrpcPort)
	log.Printf("Listening on %s\n", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("error is: %s", err.Error())
	}
	err = s.Serve(listener)
}
