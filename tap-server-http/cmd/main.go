package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
	http_tap_server "github.com/solo-io/tap-extension-examples/tap-server-http/pkg"
)

var (
	HttpPort = flag.Int("p", 8080, "port")
)

type server struct{}

func main() {
	flag.Parse()
	var dataScrubber data_scrubber.DataScrubber
	tapMessages := make(chan tap_service.TapRequest)
	dataScrubber.Init()

	address := fmt.Sprintf(":%d", *HttpPort)
	httpTapServerBuilder := http_tap_server.NewHttpTapServerBuilder().
		WithDataScrubber(dataScrubber).
		WithListenAddress(address).
		WithTapMessageChannel(tapMessages)
	tapServer := http_tap_server.NewHttpTapServer(httpTapServerBuilder)

	go tapServer.Run()
	for tapRequest := range tapMessages {
		tapRequestJson, err := json.MarshalIndent(&tapRequest, "", "  ")
		if err != nil {
			log.Printf("Error marshalling proto message to json: %s", err.Error())
		}
		log.Printf("Message contents were: %s\n", tapRequestJson)
	}
}
