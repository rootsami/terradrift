package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rootsami/terradrift/pkg/git"
	"github.com/rootsami/terradrift/pkg/tfstack"
)

// scanHandler is a handler function for scan endpoint and record metrics
// for changed resources based on the scan plan result
func (s Server) scanHandler(c *gin.Context) {

	name := c.Query("stack")
	planResp, err := tfstack.StackScan(name, s.Workdir, s.ConfigPath, s.ExtraBackendVars)

	// Reset the plan failure metric
	promMetrics.PlanFailure.With(prometheus.Labels{"stack": name}).Set(0)

	if err == nil {

		// Record metrics for drifts in resources
		promMetrics.AddResources.With(prometheus.Labels{"stack": name}).Set(float64(planResp.Add))
		promMetrics.ChangeResources.With(prometheus.Labels{"stack": name}).Set(float64(planResp.Change))
		promMetrics.DestroyResources.With(prometheus.Labels{"stack": name}).Set(float64(planResp.Destroy))

		c.JSON(200, planResp)

	} else {

		errorMessage := error.Error(err)
		if errorMessage == "stack was not found" {

			// Given stack name was not found in the configuration
			c.JSON(404, errorMessage)
		} else if strings.Contains(errorMessage, "error acquiring the state lock") {

			promMetrics.PlanFailure.With(prometheus.Labels{"stack": name}).Set(1)
			// When there's a current terrafom plan in progress, terraform locks the state till it's finished.
			c.JSON(502, "Another plan is in-progress for the requested stack, please try again in few minutes.")

		} else {

			promMetrics.PlanFailure.With(prometheus.Labels{"stack": name}).Set(1)
			c.JSON(500, errorMessage)
		}
	}
}

// gitHandler is a handler function for git sync endpoint
func (s Server) gitHandler(c *gin.Context) {

	// Reset the git failure metric
	promMetrics.GitPullFailure.With(prometheus.Labels{}).Set(0)

	status, err := git.GitPull(s.Workdir, s.GitToken, s.GitTimeout)
	if err != nil {
		promMetrics.GitPullFailure.With(prometheus.Labels{}).Set(1)
		c.JSON(500, error.Error(err))
	} else {
		c.JSON(200, status)
	}
}

// prometheusHandler returns a gin.HandlerFunc that serves prometheus metrics.
func prometheusHandler(reg *prometheus.Registry) gin.HandlerFunc {
	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
