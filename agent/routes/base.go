package routes

import (
	"github.com/gin-gonic/gin"
	agentConfig "github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/routes/config"
)

// Init the routes
func Init(g *gin.Engine, cfg agentConfig.Interface) {
	jobsRoute(g)
	statsRoute(g)
	logsRoute(g)
	config.Setup(g, cfg)
}
