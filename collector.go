package perfexporter

import (
	"runtime"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tcolgate/perfexporter/internal/system"
)

func init() {
	c := NewPerfCollector()
	prometheus.MustRegister(c)
}

type perfCollector []*system.Counter

func NewPerfCollector() prometheus.Collector {
	return newPerfCollector("")
}

func MustRegisterGoRoutineLocalCollector(name string) {
	c := newPerfCollector(name)
	prometheus.MustRegister(c)
}

func NewGoRoutineLocalCollector(name string) prometheus.Collector {
	return newPerfCollector(name)
}

func newPerfCollector(gname string) prometheus.Collector {
	pid := syscall.Getpid()
	if gname != "" {
		runtime.LockOSThread()
		pid = syscall.Gettid()
	}

	c1, _ := system.NewProcessCyclesCounter(pid)
	c2, _ := system.NewProcessInstructionCounter(pid)
	c3, _ := system.NewProcessLLCMissLoadCounter(pid)
	c4, _ := system.NewProcessLLCMissStoreCounter(pid)

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
