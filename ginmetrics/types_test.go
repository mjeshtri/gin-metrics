package ginmetrics

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSingletonRacing(t *testing.T) {
	var wg sync.WaitGroup
	nLoops := 1000
	wg.Add(nLoops)
	for i := 0; i < nLoops; i++ {
		go func() {
			GetMonitor()
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestAddMetrics(t *testing.T) {
	t.Parallel()

	m := GetMonitor()
	metric := Metric{
		Labels: nil,
	}

	err := m.AddMetric(&metric)
	require.EqualError(t, err, "metric name cannot be empty")

	// add the name
	metric.Name = metricRequestTotal
	err = m.AddMetric(&metric)
	require.EqualError(t, err, "metric type '0' is not recognized")

	// add a known metric type
	metric.Type = Counter
	err = m.AddMetric(&metric)
	require.NoError(t, err)

	err = m.AddMetric(&metric)
	require.EqualError(t, err, "metric 'gin_request_total' exists")
}
