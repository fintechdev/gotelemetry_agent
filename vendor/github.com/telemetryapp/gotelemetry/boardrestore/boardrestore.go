package main

import (
	"encoding/json"
	"flag"
	"github.com/telemetryapp/gotelemetry"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	apiKeyPtr := flag.String("api-key", "", "Telemetry API Key")
	boardNamePtr := flag.String("n", "", "Name of the board to create")
	boardPrefixPtr := flag.String("p", "", "Flow tag prefix")
	filePathPtr := flag.String("f", "", "Input file path")

	flag.Parse()

	apiKey := strings.TrimSpace(*apiKeyPtr)
	boardName := strings.TrimSpace(*boardNamePtr)
	boardPrefix := strings.TrimSpace(*boardPrefixPtr)
	filePath := strings.TrimSpace(*filePathPtr)

	if apiKey == "" {
		log.Fatal("The API Key is required")
	}

	if boardName == "" {
		log.Fatal("The board name is required")
	}

	if boardPrefix == "" {
		log.Fatal("The board prefix is required")
	}

	if filePath == "" {
		log.Fatal("The file path is required")
	}

	fileContents, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	boardTemplate := &gotelemetry.ExportedBoard{}

	err = json.Unmarshal(fileContents, &boardTemplate)

	if err != nil {
		log.Fatal(err)
	}

	credentials, err := gotelemetry.NewCredentials(apiKey)

	if err != nil {
		log.Fatal(err)
	}

	board, err := gotelemetry.ImportBoard(credentials, boardName, boardPrefix, boardTemplate)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Success: %v", board)
}
