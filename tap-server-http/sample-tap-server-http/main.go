package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/protobuf/proto"

	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
)

var (
	HttpPort     = flag.Int("p", 8080, "port")
	dataScrubber data_scrubber.DataScrubber
)

type server struct{}

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
		dataScrubber.ScrubTapRequest(tapRequest)
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
