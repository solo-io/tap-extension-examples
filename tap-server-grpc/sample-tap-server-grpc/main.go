package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	sts "sample-tap-server/tap_grpc"

	"google.golang.org/grpc"
)

var (
	GrpcPort = flag.Int("p", 8080, "port")
)

type server struct{}

func (s *server) ReportTap(srv sts.TapService_ReportTapServer) error {
	log.Printf("Starting to listen for requests")
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			log.Printf("End of stream\n")
			return ctx.Err()
		default:
		}

		req, err := srv.Recv()
		if err == io.EOF {
			// Client has closed the stream
			return nil
		}
		log.Printf("got a request!")
		req_json, err := json.Marshal(req)
		if err != nil {
			log.Printf("Unable to convert message to json, raw body is %#v\n", req.GetTraceData())
		} else {
			log.Printf("request contents are: %s\n", req_json)
		}
	}
}

func main() {
	flag.Parse()
	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(1000)}
	s := grpc.NewServer(sopts...)
	sts.RegisterTapServiceServer(s, &server{})

	address := fmt.Sprintf(":%d", *GrpcPort)
	log.Printf("Listening on %s\n", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("error is: %s", err.Error())
	}
	err = s.Serve(listener)
}
