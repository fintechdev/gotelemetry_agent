package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"log"
	"strings"
)

func main() {
	keyPtr := flag.String("api-key", "", "Telemetry API Key")
	boardNamePtr := flag.String("n", "", "Name of board to retrieve")
	boardIdPtr := flag.String("b", "", "ID of board to retrieve")

	flag.Parse()

	key := strings.TrimSpace(*keyPtr)
	boardName := strings.TrimSpace(*boardNamePtr)
	boardId := strings.TrimSpace(*boardIdPtr)

	if key == "" {
		log.Fatal("Missing API Key.")
	}

	credentials, err := gotelemetry.NewCredentials(key)

	if err != nil {
		log.Fatalf("Error reported by the Telemetry API while creating a set of credentials: %s", err.Error())
	}

	if boardName == "" && boardId == "" {
		log.Fatal("You must specify either a board ID or a board name")
	}

	if boardName != "" && boardId != "" {
		log.Fatal("You must specify *either* a board ID or a board name")
	}

	var board *gotelemetry.Board

	if boardName != "" {
		board, err = gotelemetry.GetBoardByName(credentials, boardName)

		if err != nil {
			log.Fatalf("Error reported by the Telemetry API while requesting a board by name: %s", err.Error())
		}

		boardId = board.Id
	}

	board, err = gotelemetry.GetBoard(credentials, boardId)

	if err != nil {
		log.Fatalf("Error reported by the Telemetry API while requesting a board: %s", err.Error())
	}

	exportedBoard, err := board.Export()

	if err != nil {
		log.Fatalf("Error reported by the Telemetry API while exporting a board: %s", err.Error())
	}

	result, err := json.Marshal(exportedBoard)

	if err != nil {
		log.Fatalf("Error while converting a board to JSON: %s", err.Error())
	}

	fmt.Print(string(result))
}
