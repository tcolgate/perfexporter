package perfexporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tcolgate/perfexporter/internal/system"
)

func init() {
	c := NewPerfCollector()
	prometheus.MustRegister(c)
}

type perfCollector []*system.Counter

func NewPerfCollector() prometheus.Collector {
	c1, _ := system.NewProcessCyclesCounter()
	c2, _ := system.NewProcessInstructionCounter()
	c3, _ := system.NewProcessLLCMissLoadCounter()
	c4, _ := system.NewProcessLLCMissStoreCounter()

	return perfCollector{c1, c2, c3, c4}
}

// Describe implements Collector.
func (p perfCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, c := range p {
		c.Describe(ch)
	}
}

// Collect implements Collector.
func (p perfCollector) Collect(ch chan<- prometheus.Metric) {
	for _, c := range p {
		c.Collect(ch)
	}
}
