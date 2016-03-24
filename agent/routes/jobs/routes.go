package jobs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
)

var jobManager job.JobManager

func Setup(g *gin.Engine) {
	g.GET("/jobs", Get)
	g.GET("/jobs/:id", GetByID)
	g.POST("/jobs", Post)
	g.DELETE("/jobs/:id", DeleteByID)
	jobManager = job.GetJobManager()
}

func Get(g *gin.Context) {
	for k, _ := range jobManager.Jobs {
		fmt.Println("Job ID:", k)
	}

}

func Post(g *gin.Context) {
	var jobConfig config.Job
	g.BindJSON(&jobConfig)

	err := jobManager.CreateJob(jobConfig, false)
	if err != nil {
		fmt.Println("Created Job: ", jobConfig.ID())
	}

}

func GetByID(g *gin.Context) {

	id := g.Param("id")

	if foundJob, found := jobManager.Jobs[id]; found {
		fmt.Println(foundJob)
	} else {
		fmt.Println("Job not found: ", id)
	}

}

func DeleteByID(g *gin.Context) {

	id := g.Param("id")
	job.TerminateJob(id)

}
