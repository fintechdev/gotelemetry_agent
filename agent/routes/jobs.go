package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
)

type script struct {
	Source string `json:"source" binding:"required"`
}

// jobsSetup instantiates the endpoints used for manipulating jobs
func jobsRoute(g *gin.Engine) {

	// returns a list of all jobs
	g.GET("/jobs", func(g *gin.Context) {
		job.GetJobs()
	})

	// creates a new job
	g.POST("/jobs", func(g *gin.Context) {
		var jobConfig config.Job
		g.BindJSON(&jobConfig)

		err := job.AddJob(jobConfig)
		if err != nil {
			fmt.Println("Created Job: ", jobConfig.ID)
		}
	})

	// gets a job specified by ID (returns a script if present as a nested object)
	g.GET("/jobs/:id", func(g *gin.Context) {
		id := g.Param("id")
		job.GetJobByID(id)
	})

	// deletes a job by ID
	g.DELETE("/jobs/:id", func(g *gin.Context) {
		id := g.Param("id")
		job.TerminateJob(id)
	})

	// replaces an existing job with the contents
	g.PUT("/jobs/:id", func(g *gin.Context) {
		var jobConfig config.Job
		g.BindJSON(&jobConfig)

		err := job.ReplaceJob(jobConfig)
		if err != nil {
			fmt.Println("Replaced Job: ", jobConfig.ID)
		}
	})

	// gets a script for the job
	g.GET("/jobs/:id/script", func(g *gin.Context) {
		id := g.Param("id")
		job.GetScript(id)
	})

	// creates or updates a script for a job
	g.PUT("/jobs/:id/script", func(g *gin.Context) {
		id := g.Param("id")
		var scriptSource script
		g.BindJSON(&scriptSource)
		fmt.Println(scriptSource.Source)
		job.AddScript(id, scriptSource.Source)
	})

	// removes the jobs script
	g.DELETE("/jobs/:id/script", func(g *gin.Context) {
		id := g.Param("id")
		job.DeleteScript(id)
	})

	// executes a script in debug mode and returns the output
	g.GET("/jobs/:id/script/run", func(g *gin.Context) {
		id := g.Param("id")
		job.RunScriptDebug(id)
	})

	// allow for enabling scripts
	g.PUT("/jobs/:id/script/enable", func(g *gin.Context) {
		id := g.Param("id")
		job.SetScriptState(id, true)
	})

	// allow for disabling scripts
	g.PUT("/jobs/:id/script/disable", func(g *gin.Context) {
		id := g.Param("id")
		job.SetScriptState(id, false)
	})

}
