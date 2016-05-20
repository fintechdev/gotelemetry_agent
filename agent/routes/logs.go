package routes

import "github.com/gin-gonic/gin"

func logsRoute(g *gin.Engine) {

	// returns the most recent X (say 100) log events
	g.GET("/logs", func(g *gin.Context) {
	})

	// returns an event source stream of the live log events from the agent. See EDM API for how to do eventsource streaming
	g.GET("/logs/stream", func(g *gin.Context) {
	})

}
