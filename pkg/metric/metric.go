package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var PromMetric *Metrics

type Metrics struct {
	ChangeResources  *prometheus.GaugeVec
	AddResources     *prometheus.GaugeVec
	DestroyResources *prometheus.GaugeVec
	PlanFailure      *prometheus.GaugeVec
	GitPullFailure   *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		AddResources: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "terradrift",
			Name:      "plan_add_resources",
			Help:      "Number of resources to be added based on tf plan",
		}, []string{"stack"}),
		ChangeResources: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "terradrift",
			Name:      "plan_change_resources",
			Help:      "Number of resources to be changed based on tf plan",
		}, []string{"stack"}),
		DestroyResources: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "terradrift",
			Name:      "plan_destroy_resources",
			Help:      "Number of resources to be destroyed based on tf plan",
		}, []string{"stack"}),
		PlanFailure: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "terradrift",
			Name:      "plan_failure",
			Help:      "Status of the last scan of a stack",
		}, []string{"stack"}),
		GitPullFailure: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "terradrift",
			Name:      "git_pull_failure",
			Help:      "Status of the last git pull",
		}, []string{}),
	}
	reg.MustRegister(m.AddResources)
	reg.MustRegister(m.ChangeResources)
	reg.MustRegister(m.DestroyResources)
	reg.MustRegister(m.PlanFailure)
	reg.MustRegister(m.GitPullFailure)

	return m
}
