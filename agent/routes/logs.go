package routes

import (
	"container/list"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func logsRoute(g *gin.Engine, apiStreamChannel chan string, streamRunning *bool, logList *list.List) {

	// returns the most recent X (say 100) log events
	g.GET("/logs", func(g *gin.Context) {

		var logArray []string

		for e := logList.Front(); e != nil; e = e.Next() {
			logArray = append(logArray, e.Value.(string))
		}
		g.JSON(http.StatusOK, logArray)
	})

	// returns an event source stream of the live log events from the agent
	g.GET("/logs/stream", func(g *gin.Context) {
		writer := g.Writer
		flusher, ok := writer.(http.Flusher)

		if !ok {
			// Streaming is not supported by the client
			g.AbortWithError(http.StatusBadRequest, fmt.Errorf("Streaming unsupported"))
			return
		}

		notify := writer.(http.CloseNotifier).CloseNotify()

		writer.Header().Set("Content-Type", "text/event-stream")
		writer.Header().Set("Cache-Control", "no-cache")
		writer.Header().Set("Connection", "keep-alive")

		*streamRunning = true

		for {
			select {
			case <-notify:
				*streamRunning = false
				goto Done
			case logMessage := <-apiStreamChannel:
				fmt.Fprintf(writer, "data: %s\n\n", logMessage)

				// Flush the response
				flusher.Flush()
			}
		}

	Done:
		g.Status(http.StatusNoContent)
	})

}
