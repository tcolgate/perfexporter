package system

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestIC(t *testing.T) {
	c, err := NewProcessInstructionCounter()
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
	c, err := NewProcessCyclesCounter()
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
	c1, _ := NewProcessCyclesCounter()
	c2, _ := NewProcessInstructionCounter()
	c3, _ := NewProcessLLCMissLoadCounter()
	c4, _ := NewProcessLLCMissStoreCounter()

	reg := prometheus.NewRegistry()
	reg.MustRegister(c1, c2, c3, c4)

	ms, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed, err = %v", err)
	}

	t.Logf("metrics: %s", ms)
}
