package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	agentConfig "github.com/telemetryapp/gotelemetry_agent/agent/config"
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
		g.BindJSON(&data)

		if len(data.Server.APIToken) > 0 {
			err := cfg.SetAPIToken(data.Server.APIToken)
			if err != nil {
				g.Error(err)
				return
			}
		}

		g.Status(http.StatusNoContent)
	}
}
