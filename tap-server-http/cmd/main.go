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
	HttpPort     = flag.Int("p", 8080, "port")
	OutputFormat = flag.String("output-format", "none", "which output format to use (json/none)")
)

type server struct{}
type printMessageFuncType func(*tap_service.TapRequest)

func main() {
	flag.Parse()

	// TODO deduplicate this from grpc/http main.go
	var printMessageFunc printMessageFuncType
	switch *OutputFormat {
	case "none":
		printMessageFunc = func(*tap_service.TapRequest) {}
	case "json":
		printMessageFunc = func(tapRequest *tap_service.TapRequest) {
			tapRequestJson, err := json.MarshalIndent(&tapRequest, "", "  ")
			if err != nil {
				log.Printf("Error marshalling proto message to json: %s", err.Error())
			}
			log.Printf("Message contents were: %s\n", tapRequestJson)
		}
	default:
		log.Fatalf("invalid value for --output-format: %s", *OutputFormat)
	}

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
		printMessageFunc(&tapRequest)
	}
}
