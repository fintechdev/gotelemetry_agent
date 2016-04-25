package routes

import "github.com/gin-gonic/gin"

func statsRoute(g *gin.Engine) {

	// returns stats of the agent (updates per second, error rates, etc)
	g.GET("/stats", func(g *gin.Context) {
	})

}
