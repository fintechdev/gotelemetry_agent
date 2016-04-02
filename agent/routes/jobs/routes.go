package jobs

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
)

// Setup instantiates the endpoints used for manipulating jobs
func Setup(g *gin.Engine) {
	g.GET("/jobs", get)
	g.GET("/jobs/:id", getByID)
	g.POST("/jobs", post)
	g.DELETE("/jobs/:id", deleteByID)
}

func get(g *gin.Context) {
	job.GetJobs()
}

func post(g *gin.Context) {
	var jobConfig config.Job
	g.BindJSON(&jobConfig)

	err := job.AddJob(jobConfig)
	if err != nil {
		fmt.Println("Created Job: ", jobConfig.ID)
	}
}

func getByID(g *gin.Context) {
	id := g.Param("id")
	job.GetJobByID(id)
}

func deleteByID(g *gin.Context) {
	id := g.Param("id")
	job.TerminateJob(id)
}
