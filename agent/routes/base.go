package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/telemetryapp/gotelemetry_agent/agent/routes/jobs"
)

// Init the routes
func Init(g *gin.Engine) {
	jobs.Setup(g)
}
