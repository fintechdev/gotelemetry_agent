package jobs

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
)

type script struct {
	Source string `json:"source" binding:"required"`
}

// Setup instantiates the endpoints used for manipulating jobs
func Setup(g *gin.Engine) {
	g.GET("/jobs", get)
	g.GET("/jobs/:id", getByID)
	g.POST("/jobs", post)
	g.DELETE("/jobs/:id", deleteByID)
	g.GET("/jobs/:id/script", getScript)
	g.POST("/jobs/:id/script", addScript)
	g.DELETE("/jobs/:id/script", deleteScript)
	g.GET("/jobs/:id/script/run", runScript)

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

func getScript(g *gin.Context) {
	id := g.Param("id")
	job.GetScript(id)
}

func addScript(g *gin.Context) {
	id := g.Param("id")
	var scriptSource script
	g.BindJSON(&scriptSource)
	fmt.Println(scriptSource.Source)
	job.AddScript(id, scriptSource.Source)
}

func deleteScript(g *gin.Context) {
	id := g.Param("id")
	job.DeleteScript(id)
}

func runScript(g *gin.Context) {
	id := g.Param("id")
	job.RunScriptDebug(id)
}
