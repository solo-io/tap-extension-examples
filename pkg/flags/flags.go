package flags

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
)

type printMessageFuncType func(*tap_service.TapRequest)
type Flags struct {
	Port            int
	OutputFormatter printMessageFuncType
	MessageDelay    *time.Duration
}

func ParseFlags() (*Flags, error) {
	port := flag.Int("p", 8080, "port")
	outputFormat := flag.String("output-format", "none", "which output format to use (json/none)")
	messageDelayText := `delay between acknowledgement of trace messages; can be used to simulate a slow tap server`
	messageDelay := flag.Duration("message-delay", 0, messageDelayText)
	flag.Parse()

	var printMessageFunc printMessageFuncType
	switch *outputFormat {
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
		return nil, fmt.Errorf("invalid value for --output-format: %s", *outputFormat)
	}

	return &Flags{
		Port:            *port,
		OutputFormatter: printMessageFunc,
		MessageDelay:    messageDelay,
	}, nil
}
