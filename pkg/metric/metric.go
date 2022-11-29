package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var PromMetric *Metrics

type Metrics struct {
	ChangeResources  *prometheus.GaugeVec
	AddResources     *prometheus.GaugeVec
	DestroyResources *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		AddResources: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "terradrift",
			Name:      "plan_add_resources",
			Help:      "Number of resource to be added based on tf plan",
		}, []string{"stack"}),
		ChangeResources: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "terradrift",
			Name:      "plan_change_resources",
			Help:      "Number of resource to be changed based on tf plan",
		}, []string{"stack"}),
		DestroyResources: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "terradrift",
			Name:      "plan_destroy_resources",
			Help:      "Number of resource to be destroyd based on tf plan",
		}, []string{"stack"}),
	}
	reg.MustRegister(m.AddResources)
	reg.MustRegister(m.ChangeResources)
	reg.MustRegister(m.DestroyResources)

	return m
}
