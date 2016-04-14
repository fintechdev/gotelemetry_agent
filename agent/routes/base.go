package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

// Init the routes
func Init(cfg config.Interface, errorChannel chan error) error {
	authToken := cfg.AuthToken()

	if len(authToken) == 0 {
		return nil
	}

	g := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// Handle CORS
	g.Use(func(g *gin.Context) {
		g.Header("Access-Control-Allow-Origin", "*")
		g.Header("Access-Control-Allow-Methods", "GET, PUT, PATCH, POST, DELETE, OPTIONS")
		g.Header("Access-Control-Allow-Headers", "Authorization,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")

		if g.Request.Method == "OPTIONS" {
			g.AbortWithStatus(http.StatusNoContent)
		} else {
			g.Next()
		}
	})

	// Authenticate all requests
	g.Use(func(g *gin.Context) {
		auth := g.Request.Header.Get("AUTHORIZATION")
		if strings.HasSuffix(auth, authToken) {
			g.Next()
		} else {
			g.AbortWithStatus(http.StatusUnauthorized)
		}
	})

	// Handle errors
	g.Use(func(g *gin.Context) {
		g.Next()

		if len(g.Errors) == 0 {
			return
		}

		var errMesages []string
		status := 500

		for _, err := range g.Errors {
			errMesages = append(errMesages, err.Error())

			if err.IsType(gin.ErrorTypePublic) {
				status = 520
			} else if err.IsType(gin.ErrorTypeRender) {
				status = 444
			} else if err.IsType(gin.ErrorTypeBind) {
				status = 400
			} else if err.IsType(gin.ErrorTypeNu) {
				status = 521
			} else if err.IsType(gin.ErrorTypeAny) {
				status = 522
			}

			errorChannel <- gotelemetry.NewError(status, err.Error())
		}

		g.JSON(status, gin.H{
			"code":   status,
			"errors": errMesages})
	})

	jobsRoute(g)
	statsRoute(g)
	logsRoute(g)
	configRoute(g, cfg)

	listen := cfg.Listen()
	if len(listen) == 0 {
		listen = ":8080"
	}
	errorChannel <- gotelemetry.NewLogError("Listening at %s", listen)
	go g.Run(listen)

	return nil
}
