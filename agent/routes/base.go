package routes

import (
	"github.com/gin-gonic/gin"
	agentConfig "github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/routes/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/routes/jobs"
)

// Init the routes
func Init(g *gin.Engine, cfg agentConfig.Interface) {
	jobs.Setup(g)
	config.Setup(g, cfg)
}
