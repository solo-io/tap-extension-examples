package tap_server_builder

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type TapServer interface {
	Run(listenAddress string)
}

type httpTapServerImpl struct {
	handler func(http.ResponseWriter, *http.Request)
}

type grpcTapServerImpl struct {
	tapMessages    chan tap_service.TapRequest
	dataScrubber   data_scrubber.DataScrubber
	grpcServerOpts []grpc.ServerOption
}

func (tapServerImpl *grpcTapServerImpl) ReportTap(srv tap_service.TapService_ReportTapServer) error {
	log.Printf("Starting to listen for tap reports")
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
		log.Printf("got a tap report!")
		tapServerImpl.dataScrubber.ScrubTapRequest(tapRequest)
		tapServerImpl.tapMessages <- *tapRequest
	}
}

type tapServerBuilder struct {
	// channel where received tap requests will be written
	tapMessages chan tap_service.TapRequest
	// data scrubber - can be run on TapRequest objects to purge sensitive data
	// prior to being written on tapMessages
	dataScrubber data_scrubber.DataScrubber
}

func NewTapServerBuilder() *tapServerBuilder {
	return &tapServerBuilder{}
}

func (tapServerBuilder *tapServerBuilder) WithTapMessageChannel(tapMessages chan tap_service.TapRequest) *tapServerBuilder {
	tapServerBuilder.tapMessages = tapMessages
	return tapServerBuilder
}

func (tapServerBuilder *tapServerBuilder) WithDataScrubber(dataScrubber data_scrubber.DataScrubber) *tapServerBuilder {
	tapServerBuilder.dataScrubber = dataScrubber
	return tapServerBuilder
}

func (tapServerBuilder *tapServerBuilder) BuildHttp() TapServer {
	handler := func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("got a request on %s\n", r.URL.Path)
		traceData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request traceData: %s", err.Error())
			return
		}
		tapRequest := &tap_service.TapRequest{}
		proto.Unmarshal(traceData, tapRequest)
		tapServerBuilder.dataScrubber.ScrubTapRequest(tapRequest)
		tapServerBuilder.tapMessages <- *tapRequest
	}
	return &httpTapServerImpl{
		handler: handler,
	}
}

func (tapServerBuilder *tapServerBuilder) BuildGrpc(grpcServerOpts []grpc.ServerOption) TapServer {
	return &grpcTapServerImpl{
		tapMessages:    tapServerBuilder.tapMessages,
		dataScrubber:   tapServerBuilder.dataScrubber,
		grpcServerOpts: grpcServerOpts,
	}
}

func (tap_server *httpTapServerImpl) Run(listenAddress string) {
	log.Printf("Listening on %s\n", listenAddress)
	err := http.ListenAndServe(listenAddress, http.HandlerFunc(tap_server.handler))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (tap_server *grpcTapServerImpl) Run(listenAddress string) {
	log.Printf("Listening on %s\n", listenAddress)
	server := grpc.NewServer(tap_server.grpcServerOpts...)
	tap_service.RegisterTapServiceServer(server, tap_server)
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal("error starting server", err)
	}
	log.Fatal(server.Serve(listener))
}
