package jobs

import (
	"github.com/gin-gonic/gin"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
)

// Setup the routes
func Setup(g *gin.Engine) {
	g.GET("/jobs", Get)
}

// Add a new Page
func Get(g *gin.Context) {
	job.GetJobsList()
}
