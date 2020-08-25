package monitoring

import (
	"time"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

var (
	ms = time.Millisecond.Seconds()

	// defaultHistogramBucketsForTimer is a default histogram buckets, used mostly for Timer.
	//
	// They are mapped after https://github.com/prometheus/client_golang/blob/master/prometheus/histogram.go#L62
	defaultHistogramBucketsForTimer monitor.Buckets = []float64{
		5 * ms,
		10 * ms,
		25 * ms,
		50 * ms,
		100 * ms,
		250 * ms,
		500 * ms,
		1000 * ms,
		2500 * ms,
		5000 * ms,
		10000 * ms,
	}
)

type mortarMetric struct {
	*tagsMetric
	registry *externalRegistry
	cfg      *monitorConfig
}

func newMetric(externalMonitor monitor.BricksMetrics, cfg *monitorConfig) monitor.Metrics {
	return &mortarMetric{
		registry:   newRegistry(externalMonitor),
		cfg:        cfg,
		tagsMetric: &tagsMetric{tags: cfg.tags},
	}
}

// Counter creates a counter with possible predefined tags
func (mm *mortarMetric) Counter(name, desc string) monitor.TagsAwareCounter {
	bricksCounter, err := mm.registry.loadOrStoreCounter(name, desc, mm.extractTagKeys()...)
	if err != nil {
		mm.cfg.onError(err)
		bricksCounter = newNoopCounter(err, mm.cfg.onError)
	}

	return newCounterWithTags(bricksCounter, mm.tags, mm.cfg.extractors)
}

// Gauge creates a gauge with possible predefined tags
func (mm *mortarMetric) Gauge(name, desc string) monitor.TagsAwareGauge {
	bricksGauge, err := mm.registry.loadOrStoreGauge(name, desc, mm.extractTagKeys()...)
	if err != nil {
		mm.cfg.onError(err)
		bricksGauge = newNoopGauge(err, mm.cfg.onError)
	}

	return newGaugeWithTags(bricksGauge, mm.tags, mm.cfg.extractors)
}

// Histogram creates a histogram with possible predefined tags
func (mm *mortarMetric) Histogram(name, desc string, buckets monitor.Buckets) monitor.TagsAwareHistogram {
	bricksHistogram, err := mm.registry.loadOrStoreHistogram(name, desc, buckets, mm.extractTagKeys()...)
	if err != nil {
		mm.cfg.onError(err)
		bricksHistogram = newNoopHistogram(err, mm.cfg.onError)
	}

	return newHistogramWithTags(bricksHistogram, mm.tags, mm.cfg.extractors)
}

// Timer creates a timer with possible predefined tags
func (mm *mortarMetric) Timer(name, desc string, buckets monitor.Buckets) monitor.TagsAwareTimer {
	if buckets == nil {
		buckets = defaultHistogramBucketsForTimer
	}
	bricksHistogram, err := mm.registry.loadOrStoreHistogram(name, desc, buckets, mm.extractTagKeys()...)
	if err != nil {
		mm.cfg.onError(err)
		bricksHistogram = newNoopHistogram(err, mm.cfg.onError)
	}

	return newTimerWithTags(bricksHistogram, mm.tags, mm.cfg.extractors)
}

// WithTags sets custom tags to be included if possible in every Metric
func (mm *mortarMetric) WithTags(tags monitor.Tags) monitor.Metrics {
	mm.withTags(tags)
	return mm
}

func (mm *mortarMetric) extractTagKeys() (keys []string) {
	for k := range mm.tags {
		keys = append(keys, k)
	}
	return
}