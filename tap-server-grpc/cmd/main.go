package main

import (
	"fmt"
	"log"

	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	"github.com/solo-io/tap-extension-examples/pkg/flags"
	"github.com/solo-io/tap-extension-examples/pkg/tap_server_builder"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
	"google.golang.org/grpc"
)

func main() {
	config, err := flags.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	var dataScrubber data_scrubber.DataScrubber
	dataScrubber.Init()
	tapMessages := make(chan tap_service.TapRequest)

	listenAddress := fmt.Sprintf(":%d", config.Port)
	tapServerBuilder := tap_server_builder.NewTapServerBuilder().
		WithDataScrubber(&dataScrubber).
		WithTapMessageChannel(tapMessages)
	tapServer := tapServerBuilder.BuildGrpc([]grpc.ServerOption{grpc.MaxConcurrentStreams(1000)})

	go tapServer.Run(listenAddress)
	for tapRequest := range tapMessages {
		config.OutputFormatter(&tapRequest)
	}
}
