package routes

import (
"net/http"
"github.com/gin-gonic/gin"
)


func logsRoute(g *gin.Engine, apiStreamChannel chan string) {

	// returns the most recent X (say 100) log events
	g.GET("/logs", func(g *gin.Context) {
	})

	// returns an event source stream of the live log events from the agent. See EDM API for how to do eventsource streaming
	g.GET("/logs/stream", func(g *gin.Context) {

		g.Header("Content-Type", "text/event-stream")
		g.Header("Cache-Control", "no-cache")
		g.Header("Connection", "keep-alive")
		g.Status(http.StatusNoContent)

			for {
				select {
				case logMessage := <-apiStreamChannel:
					g.SSEvent("log", logMessage)
				}
			}
	})

}
