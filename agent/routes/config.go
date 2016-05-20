package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	agentConfig "github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/database"
)

type request struct {
	Server   server   `json:"server"`
	Data     data     `json:"data"`
	Graphite graphite `json:"graphite"`
}

type server struct {
	APIToken string `json:"api_token"`
}

type data struct {
	TTL string `json:"ttl"`
}

type graphite struct {
	TCPListenPort string `json:"listen_tcp"`
	UDPListenPort string `json:"listen_udp"`
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

		if len(data.Data.TTL) > 0 {
			database.WriteConfigParam("ttl", data.Data.TTL)
			cfg.SetDatabaseTTL(data.Data.TTL)
		}

		if len(data.Graphite.UDPListenPort) > 0 {
			database.WriteConfigParam("listen_udp", data.Graphite.UDPListenPort)
			cfg.SetUDPListenPort(data.Graphite.UDPListenPort)
		} else if len(data.Graphite.TCPListenPort) > 0 {
			database.WriteConfigParam("listen_tcp", data.Graphite.TCPListenPort)
			cfg.SetTCPListenPort(data.Graphite.TCPListenPort)
		}

		g.Status(http.StatusNoContent)
	}
}
