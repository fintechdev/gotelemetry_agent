package jobs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

func Setup(g *gin.Engine) {
	g.GET("/jobs", Get)
	g.POST("/jobs", Post)
}

func Get(g *gin.Context) {
	job.GetJobsList()
}

type Job config.Job

// TODO finish POST function
func Post(g *gin.Context) {
	var job Job
	g.BindJSON(&job)
	fmt.Println(job)
}
