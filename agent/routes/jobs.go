package routes

import (
	"net/http"
	"net/url"

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
		jobsList, _ := job.GetJobs()
		g.JSON(http.StatusOK, jobsList)
	})

	// creates a new job
	g.POST("/jobs", func(g *gin.Context) {
		var jobConfig config.Job
		if err := g.BindJSON(&jobConfig); err != nil {
			g.Error(err).SetType(gin.ErrorTypeBind)
			return
		}

		err := job.AddJob(jobConfig)
		if err != nil {
			g.Error(err)
			return
		}

		g.Status(http.StatusNoContent)
	})

	// gets a job specified by ID (returns a script if present as a nested object)
	g.GET("/jobs/:id", func(g *gin.Context) {
		id, _ := url.QueryUnescape(g.Param("id"))
		jobConfig, err := job.GetJobByID(id)

		if err != nil {
			g.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "errors": err.Error()})
			return
		}

		g.JSON(http.StatusOK, jobConfig)
	})

	// deletes a job by ID
	g.DELETE("/jobs/:id", func(g *gin.Context) {
		id, _ := url.QueryUnescape(g.Param("id"))
		err := job.TerminateJob(id)

		if err != nil {
			g.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "errors": err.Error()})
			return
		}

		g.Status(http.StatusNoContent)
	})

	// replaces an existing job with the contents
	g.PATCH("/jobs/:id", func(g *gin.Context) {
		var jobConfig config.Job
		if err := g.BindJSON(&jobConfig); err != nil {
			g.Error(err).SetType(gin.ErrorTypeBind)
			return
		}

		id, _ := url.QueryUnescape(g.Param("id"))
		jobConfig.ID = id

		err := job.ReplaceJob(jobConfig)

		if err != nil {
			g.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "errors": err.Error()})
			return
		}

		g.Status(http.StatusNoContent)
	})

	// gets a script for the job
	g.GET("/jobs/:id/script", func(g *gin.Context) {
		id, _ := url.QueryUnescape(g.Param("id"))
		jobScript, err := job.GetScript(id)

		if err != nil {
			g.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "errors": err.Error()})
			return
		}

		jobScriptMap := make(map[string]string)
		jobScriptMap["source"] = jobScript

		g.JSON(http.StatusOK, jobScriptMap)
	})

	// creates or updates a script for a job
	g.POST("/jobs/:id/script", func(g *gin.Context) {
		id, _ := url.QueryUnescape(g.Param("id"))

		var scriptSource script
		if err := g.BindJSON(&scriptSource); err != nil {
			g.Error(err).SetType(gin.ErrorTypeBind)
			return
		}

		err := job.AddScript(id, scriptSource.Source)

		if err != nil {
			g.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "errors": err.Error()})
			return
		}

		g.Status(http.StatusNoContent)
	})

	// removes the jobs script
	g.DELETE("/jobs/:id/script", func(g *gin.Context) {
		id, _ := url.QueryUnescape(g.Param("id"))
		err := job.DeleteScript(id)

		if err != nil {
			g.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "errors": err.Error()})
			return
		}

		g.Status(http.StatusNoContent)
	})

	// executes a script in debug mode and returns the output
	g.GET("/jobs/:id/script/run", func(g *gin.Context) {
		id, _ := url.QueryUnescape(g.Param("id"))
		res, err := job.RunScriptDebug(id)

		if err != nil {
			g.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "errors": err.Error()})
			return
		}

		g.JSON(http.StatusOK, res)
	})

	// allow for enabling scripts
	g.POST("/jobs/:id/script/enable", func(g *gin.Context) {
		id, _ := url.QueryUnescape(g.Param("id"))
		err := job.SetScriptState(id, true)

		if err != nil {
			g.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "errors": err.Error()})
			return
		}

		g.Status(http.StatusNoContent)
	})

	// allow for disabling scripts
	g.POST("/jobs/:id/script/disable", func(g *gin.Context) {
		id, _ := url.QueryUnescape(g.Param("id"))
		err := job.SetScriptState(id, false)

		if err != nil {
			g.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "errors": err.Error()})
			return
		}

		g.Status(http.StatusNoContent)
	})

}
