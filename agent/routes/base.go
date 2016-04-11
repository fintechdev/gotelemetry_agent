package routes

import "github.com/gin-gonic/gin"

// Init the routes
func Init(g *gin.Engine) {
	jobsRoute(g)
	statsRoute(g)
	logsRoute(g)
}
