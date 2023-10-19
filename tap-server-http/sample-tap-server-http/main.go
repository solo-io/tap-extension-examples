package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	tapv3 "github.com/envoyproxy/go-control-plane/envoy/data/tap/v3"
	"google.golang.org/protobuf/proto"

	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
)

var (
	HttpPort     = flag.Int("p", 8080, "port")
	dataScrubber data_scrubber.DataScrubber
)

type server struct{}

func scrubTapRequest(tapRequest *tap_service.TapRequest) {
	scrubHeader := func(header *corev3.HeaderValue) {
		fmt.Printf("\theaders are: %s:%s\n", header.GetKey(), header.GetValue())
		header.Value = dataScrubber.ScrubDataString(header.Value)
	}
	scrubBody := func(body *tapv3.Body) {
		dataScrubber.ScrubData(body.GetAsBytes())
	}
	var trace *tapv3.HttpBufferedTrace
	trace = tapRequest.GetTraceData().GetHttpBufferedTrace()
	request := trace.GetRequest()
	response := trace.GetResponse()
	fmt.Printf("Parsing request headers\n")
	for _, header := range request.GetHeaders() {
		scrubHeader(header)
	}
	fmt.Printf("Parsing response headers\n")
	for _, header := range response.GetHeaders() {
		scrubHeader(header)
	}
	scrubBody(request.GetBody())
	scrubBody(response.GetBody())
}

func main() {
	flag.Parse()
	dataScrubber.Init()

	handler := func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("got a request on %s\n", r.URL.Path)
		traceData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request traceData: %s", err.Error())
			return
		}
		tapRequest := &tap_service.TapRequest{}
		proto.Unmarshal(traceData, tapRequest)
		scrubTapRequest(tapRequest)
		tapRequestJson, err := json.MarshalIndent(tapRequest, "", "  ")
		if err != nil {
			log.Printf("Error marshalling proto message to json: %s", err.Error())
		}
		log.Printf("Message contents were: %s\n", tapRequestJson)
	}

	address := fmt.Sprintf(":%d", *HttpPort)
	log.Printf("Listening on %s\n", address)
	err := http.ListenAndServe(address, http.HandlerFunc(handler))
	if err != nil {
		log.Fatal(err.Error())
	}
}
