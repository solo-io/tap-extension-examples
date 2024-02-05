package tap_server_builder

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type TapServer interface {
	Run(listenAddress string)
	Stop() error
}

type httpTapServerImpl struct {
	handler    func(http.ResponseWriter, *http.Request)
	httpServer http.Server
}

type grpcTapServerImpl struct {
	tapMessages    chan tap_service.TapRequest
	dataScrubber   *data_scrubber.DataScrubber
	messageDelay   *time.Duration
	grpcServerOpts []grpc.ServerOption
	grpcServer     *grpc.Server
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
		if tapServerImpl.dataScrubber != nil {
			tapServerImpl.dataScrubber.ScrubTapRequest(tapRequest)
		}
		time.Sleep(*tapServerImpl.messageDelay)
		tapServerImpl.tapMessages <- *tapRequest
	}
}

type tapServerBuilder struct {
	// channel where received tap requests will be written
	tapMessages chan tap_service.TapRequest
	// data scrubber - can be run on TapRequest objects to purge sensitive data
	// prior to being written on tapMessages. can be set to nil to disable this
	// functionality
	dataScrubber *data_scrubber.DataScrubber
	// delay between acknowledgement of trace messages
	messageDelay *time.Duration
}

func NewTapServerBuilder() *tapServerBuilder {
	return &tapServerBuilder{}
}

func (tapServerBuilder *tapServerBuilder) WithTapMessageChannel(tapMessages chan tap_service.TapRequest) *tapServerBuilder {
	tapServerBuilder.tapMessages = tapMessages
	return tapServerBuilder
}

func (tapServerBuilder *tapServerBuilder) WithDataScrubber(dataScrubber *data_scrubber.DataScrubber) *tapServerBuilder {
	tapServerBuilder.dataScrubber = dataScrubber
	return tapServerBuilder
}

func (tapServerBuilder *tapServerBuilder) WithMessageDelay(messageDelay *time.Duration) *tapServerBuilder {
	tapServerBuilder.messageDelay = messageDelay
	return tapServerBuilder
}

func (tapServerBuilder *tapServerBuilder) BuildHttp() TapServer {
	handler := func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("got a request on %s\n", r.URL.Path)
		traceData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request traceData: %s", err.Error())
			return
		}
		tapRequest := &tap_service.TapRequest{}
		proto.Unmarshal(traceData, tapRequest)
		if tapServerBuilder.dataScrubber != nil {
			tapServerBuilder.dataScrubber.ScrubTapRequest(tapRequest)
		}
		tapServerBuilder.tapMessages <- *tapRequest
		if tapServerBuilder.messageDelay != nil {
			time.Sleep(*tapServerBuilder.messageDelay)
		}
	}
	return &httpTapServerImpl{
		handler: handler,
	}
}

func (tapServerBuilder *tapServerBuilder) BuildGrpc(grpcServerOpts []grpc.ServerOption) TapServer {
	return &grpcTapServerImpl{
		tapMessages:    tapServerBuilder.tapMessages,
		dataScrubber:   tapServerBuilder.dataScrubber,
		messageDelay:   tapServerBuilder.messageDelay,
		grpcServerOpts: grpcServerOpts,
	}
}

func (tap_server *httpTapServerImpl) Run(listenAddress string) {
	log.Printf("Listening on %s\n", listenAddress)
	tap_server.httpServer.Addr = listenAddress
	tap_server.httpServer.Handler = http.HandlerFunc(tap_server.handler)
	if err := tap_server.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}
}

func (tap_server *httpTapServerImpl) Stop() error {
	if tap_server == nil {
		return nil
	}
	return tap_server.httpServer.Shutdown(context.Background())
}

func (tap_server *grpcTapServerImpl) Run(listenAddress string) {
	log.Printf("Listening on %s\n", listenAddress)
	tap_server.grpcServer = grpc.NewServer(tap_server.grpcServerOpts...)
	tap_service.RegisterTapServiceServer(tap_server.grpcServer, tap_server)
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal("error starting server", err)
	}
	log.Fatal(tap_server.grpcServer.Serve(listener))
}

func (tap_server *grpcTapServerImpl) Stop() error {
	if tap_server == nil {
		return nil
	}
	if tap_server.grpcServer == nil {
		return nil
	}
	tap_server.grpcServer.Stop()
	return nil
}
