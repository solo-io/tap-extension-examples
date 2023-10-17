package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"

	sts "sample-tap-server/tap_grpc"

	"google.golang.org/grpc"
)

var (
	GrpcPort = flag.String("port", ":9001", "port")
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
	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(1000)}
	s := grpc.NewServer(sopts...)
	sts.RegisterTapServiceServer(s, &server{})

	log.Printf("Listening on %s\n", *GrpcPort)
	listener, err := net.Listen("tcp", *GrpcPort)
	if err != nil {
		log.Printf("error is: %s", err.Error())
	}
	err = s.Serve(listener)
}
