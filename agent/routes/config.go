package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	agentConfig "github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/database"
)

type request struct {
	Server struct {
		APIToken string `json:"api_token"`
	} `json:"server"`
}

// Setup instantiates the endpoints used for manipulating jobs
func configRoute(g *gin.Engine, cfg agentConfig.Interface) {
	g.PATCH("/config", updateFunc(cfg))
}

func updateFunc(cfg agentConfig.Interface) gin.HandlerFunc {
	return func(g *gin.Context) {
		data := request{}
		if err := g.BindJSON(&data); err != nil {
			g.Error(err).SetType(gin.ErrorTypeBind)
			return
		}

		if len(data.Server.APIToken) > 0 {
			database.WriteConfigParam("api_token", data.Server.APIToken)
			cfg.SetAPIToken(data.Server.APIToken)
		}

		g.Status(http.StatusNoContent)
	}
}
