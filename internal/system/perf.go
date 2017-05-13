package system

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

func NewProcessInstructionCounter(pid int, gname string) (*Counter, error) {
	name := "linux_perf_cpu_instructions_total"
	var labels prometheus.Labels
	if gname != "" {
		name = "linux_perf_goroutine_cpu_instructions_total"
		labels = prometheus.Labels{"goroutine": gname}
	}
	help := "count of instructions despatched since process start"

	fd, err := newProcessEventCounter(perfTypeHardware, perfCountHWInstructions, pid)
	if err != nil {
		return nil, err
	}
	return &Counter{
		fd:   fd,
		desc: prometheus.NewDesc(name, help, nil, labels),
	}, nil
}

func NewProcessCyclesCounter(pid int, gname string) (*Counter, error) {
	name := "linux_perf_cpu_cycles_total"
	var labels prometheus.Labels
	if gname != "" {
		name = "linux_perf_goroutine_cpu_cycles_total"
		labels = prometheus.Labels{"goroutine": gname}
	}
	help := "count of cpu cycles since process start"

	fd, err := newProcessEventCounter(perfTypeHardware, perfCountHWCPUCycles, pid)
	if err != nil {
		return nil, err
	}
	return &Counter{
		fd:   fd,
		desc: prometheus.NewDesc(name, help, nil, labels),
	}, nil
}

func NewProcessLLCMissLoadCounter(pid int, gname string) (*Counter, error) {
	name := "linux_perf_cache_miss_total"
	labels := prometheus.Labels{"op": "load"}
	if gname != "" {
		name = "linux_perf_goroutine_cache_miss_total"
		labels["goroutine"] = gname
	}
	help := "count of cache misses since process start"

	fd, err := newProcessEventCounter(
		perfTypeHHCache,
		perfCountHWCacheLL*perfCountHWCacheResultMiss*perfCountHWCacheOpRead, pid)
	if err != nil {
		return nil, err
	}
	return &Counter{
		fd:   fd,
		desc: prometheus.NewDesc(name, help, nil, labels),
	}, nil
}

func NewProcessLLCMissStoreCounter(pid int, gname string) (*Counter, error) {
	name := "linux_perf_cache_miss_total"
	labels := prometheus.Labels{"op": "store"}
	if gname != "" {
		name = "linux_perf_goroutine_cache_miss_total"
		labels["goroutine"] = gname
	}
	help := "count of cache misses since process start"

	fd, err := newProcessEventCounter(
		perfTypeHHCache,
		perfCountHWCacheLL*perfCountHWCacheResultMiss*perfCountHWCacheOpWrite, pid)
	if err != nil {
		return nil, err
	}
	return &Counter{
		fd:   fd,
		desc: prometheus.NewDesc(name, help, nil, labels),
	}, nil
}

type Counter struct {
	name   string
	help   string
	fd     *os.File
	labels prometheus.Labels
	values []string
	desc   *prometheus.Desc
}

func (c *Counter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *Counter) Collect(ch chan<- prometheus.Metric) {
	v, _ := c.read()
	ch <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, v, c.values...)
}

const (
	perfTypeHardware   = 0
	perfTypeSoftware   = 1
	perfTypeTracepoint = 2
	perfTypeHHCache    = 3
	perfTypeRaw        = 4
	perfTypeBreakpoint = 5

	perfCountHWCPUCycles             = 0
	perfCountHWInstructions          = 1
	perfCountHWCacheReferences       = 2
	perfCountHWCacheMisses           = 3
	perfCountHWBranchInstructions    = 4
	perfCountHWBranchMisses          = 5
	perfCountHWBusCycles             = 6
	perfCountHWStalledCyclesFrontend = 7
	perfCountHWStalledCyclesBackend  = 8
	perfCountHWRefCPUCycles          = 9

	/*
	   Generalized hardware cache events:

	         { L1-D, L1-I, LLC, ITLB, DTLB, BPU, NODE } x
	         { read, write, prefetch } x
	         { accesses, misses }
	*/
	perfCountHWCacheL1d  = 0
	perfCountHWCacheL1i  = 1
	perfCountHWCacheLL   = 2
	perfCountHWCacheDTLB = 3
	perfCountHWCacheITLB = 4
	perfCountHWCacheBPU  = 5
	perfCountHWCacheNode = 6

	perfCountHWCacheOpRead     = 0
	perfCountHWCacheOpWrite    = 1
	perfCountHWCacheOpPrefetch = 2

	perfCountHWCacheResultAccess = 0
	perfCountHWCacheResultMiss   = 1

	perfCountSWCPUClock        = 0
	perfCountSWTaskClock       = 1
	perfCountSWPageFaults      = 2
	perfCountSWContextSwitches = 3
	perfCountSWCPUMigrations   = 4
	perfCountSWPageFaultsMin   = 5
	perfCountSWPageFaultsMaj   = 6
	perfCountSWAlignmentFaults = 7
	perfCountSWEmulationFaults = 8
	perfCountSWDummy           = 9
	perfCountSWBBPFOutput      = 10

	perfRecordMiscCPUModeUnknown = 0
	perfRecordMiscKernel         = 1
	perfRecordMiscUser           = 2
	perfRecordMiscHypervisor     = 3
	perfRecordMiscGuestKernel    = 4
	perfRecordMiscGuestUser      = 5
	perfRecordMiscCPUModeMask    = 7
)

const (
	perfFlagDisabled = 1 << iota
	perfFlagInherit
	perfFlagPinned
	perfFlagExclusive
	perfFlagExcludeUser
	perfFlagExcludeKernel
	perfFlagExcludeHv
	perfFlagExcludeIdle
	perfFlagMmap
	perfFlagComm
	perfFlagFreq
	perfFlagInheritStat
	perfFlagEnableOnExec
	perfFlagTask
	perfFlagWatermark
	perfFlagPreciseIPLow
	perfFlagPreciseIPHigh
	perfFlagMmapData
	perfFlagSampleIDAll
	perfFlagExcludeHost
	perfFlagExcludeGuest
	perfFlagExcludeCallchainKernel
	perfFlagExcludeCallchainUser
	perfFlagMmap2
	perfFlagCommExec
	perfFlagUseClockID
	perfFlagContextSwitch
	perfFlagWriteBackward
)
