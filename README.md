# perf-events prom collection

This is a PoC of collecting linux perf-events
counters via prometheus.

At the moment this just adds for hard coded stats based on perf
counters

-
-
-
-

## Usage

```
package main

import (
        "net/http"

        "github.com/prometheus/client_golang/prometheus/promhttp"
        _ "github.com/tcolgate/perfexporter"
)

func main() {
        // ...
        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":8011", nil)
}

```

## Example Outpu

```
# HELP linux_perf_cache_miss_total count of cache misses since process start
# TYPE linux_perf_cache_miss_total counter
linux_perf_cache_miss_total{op="load"} 1.9986703802e+10
linux_perf_cache_miss_total{op="store"} 1.2278385611e+10
# HELP linux_perf_cpu_cycles_total count of cpu cycles since process start
# TYPE linux_perf_cpu_cycles_total counter
linux_perf_cpu_cycles_total 3.52313006355e+11
# HELP linux_perf_cpu_instructions_total count of instructions despatched since process start
# TYPE linux_perf_cpu_instructions_total counter
linux_perf_cpu_instructions_total 1.67387156785e+11
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
```

