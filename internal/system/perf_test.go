package system

import (
	"syscall"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestIC(t *testing.T) {
	pid := syscall.Getpid()
	c, err := NewProcessInstructionCounter(pid, "")
	if err != nil {
		t.Fatalf("creating counter failed, err = %v", err)
	}

	v, err := c.read()
	if err != nil {
		t.Fatalf("reading counter failed, err = %v", err)
	}
	t.Logf("v = %v", v)
}

func TestCC(t *testing.T) {
	pid := syscall.Getpid()
	c, err := NewProcessCyclesCounter(pid, "")
	if err != nil {
		t.Fatalf("creating counter failed, err = %v", err)
	}

	v, err := c.read()
	if err != nil {
		t.Fatalf("reading counter failed, err = %v", err)
	}

	t.Logf("v = %v", v)
}

func TestCollector(t *testing.T) {
	pid := syscall.Getpid()
	c1, _ := NewProcessCyclesCounter(pid, "")
	c2, _ := NewProcessInstructionCounter(pid, "")
	c3, _ := NewProcessLLCMissLoadCounter(pid, "")
	c4, _ := NewProcessLLCMissStoreCounter(pid, "")

	reg := prometheus.NewRegistry()
	reg.MustRegister(c1, c2, c3, c4)

	ms, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed, err = %v", err)
	}

	t.Logf("metrics: %s", ms)
}
