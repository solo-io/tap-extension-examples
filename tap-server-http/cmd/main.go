package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	"github.com/solo-io/tap-extension-examples/pkg/tap_server_builder"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
)

var (
	HttpPort = flag.Int("p", 8080, "port")
)

type server struct{}

func main() {
	flag.Parse()
	var dataScrubber data_scrubber.DataScrubber
	dataScrubber.Init()
	tapMessages := make(chan tap_service.TapRequest)

	listenAddress := fmt.Sprintf(":%d", *HttpPort)
	httpTapServerBuilder := tap_server_builder.NewTapServerBuilder().
		WithDataScrubber(&dataScrubber).
		WithTapMessageChannel(tapMessages)
	tapServer := httpTapServerBuilder.BuildHttp()

	go tapServer.Run(listenAddress)
	for tapRequest := range tapMessages {
		tapRequestJson, err := json.MarshalIndent(&tapRequest, "", "  ")
		if err != nil {
			log.Printf("Error marshalling proto message to json: %s", err.Error())
		}
		log.Printf("Message contents were: %s\n", tapRequestJson)
	}
}
