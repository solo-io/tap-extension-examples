package http_tap_server

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
	"google.golang.org/protobuf/proto"
)

type HttpTapServer interface {
	Run()
}

type httpTapServerImpl struct {
	handler       func(http.ResponseWriter, *http.Request)
	listenAddress string
}

type httpTapServerBuilder struct {
	tapMessages   chan tap_service.TapRequest
	dataScrubber  data_scrubber.DataScrubber
	listenAddress string
}

func NewHttpTapServerBuilder() *httpTapServerBuilder {
	return &httpTapServerBuilder{}
}

func (tap_server_builder *httpTapServerBuilder) WithTapMessageChannel(tapMessages chan tap_service.TapRequest) *httpTapServerBuilder {
	tap_server_builder.tapMessages = tapMessages
	return tap_server_builder
}

func (tap_server_builder *httpTapServerBuilder) WithDataScrubber(dataScrubber data_scrubber.DataScrubber) *httpTapServerBuilder {
	tap_server_builder.dataScrubber = dataScrubber
	return tap_server_builder
}

func (tap_server_builder *httpTapServerBuilder) WithListenAddress(listenAddress string) *httpTapServerBuilder {
	tap_server_builder.listenAddress = listenAddress
	return tap_server_builder
}

func NewHttpTapServer(tapServerBuilder *httpTapServerBuilder) HttpTapServer {
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
		handler:       handler,
		listenAddress: tapServerBuilder.listenAddress,
	}
}

func (tap_server *httpTapServerImpl) Run() {
	log.Printf("Listening on %s\n", tap_server.listenAddress)
	err := http.ListenAndServe(tap_server.listenAddress, http.HandlerFunc(tap_server.handler))
	if err != nil {
		log.Fatal(err.Error())
	}
}
