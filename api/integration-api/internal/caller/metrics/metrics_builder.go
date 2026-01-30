package internal_caller_metrics

import (
	"fmt"
	"time"

	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/protos"
)

type MetricBuilder struct {
	metricsMap map[string]*protos.Metric
	start      time.Time
	requestId  uint64
}

// NewMetricBuilder initializes and returns a new MetricBuilder
func NewMetricBuilder(requestId uint64) *MetricBuilder {
	return &MetricBuilder{
		metricsMap: make(map[string]*protos.Metric),
		requestId:  requestId,
	}
}

// OnStart starts the timer and initializes basic metrics
func (mb *MetricBuilder) OnStart() *MetricBuilder {
	mb.start = time.Now()

	mb.metricsMap[type_enums.TIME_TAKEN.String()] = &protos.Metric{
		Name:        type_enums.TIME_TAKEN.String(),
		Value:       fmt.Sprintf("%d", int64(time.Since(mb.start))),
		Description: "Time taken to serve the llm request",
	}

	mb.metricsMap[type_enums.LLM_REQUEST_ID.String()] = &protos.Metric{
		Name:        type_enums.LLM_REQUEST_ID.String(),
		Value:       fmt.Sprintf("%d", mb.requestId),
		Description: "LLM Request ID",
	}

	mb.metricsMap[type_enums.STATUS.String()] = &protos.Metric{
		Name:        type_enums.STATUS.String(),
		Value:       type_enums.RECORD_FAILED.String(), // Initially mark as RECORD_FAILED
		Description: "Status of the given request to LLM",
	}

	return mb
}

// OnSuccess updates the time taken and status metrics for success
func (mb *MetricBuilder) OnSuccess() *MetricBuilder {
	if metric, exists := mb.metricsMap[type_enums.TIME_TAKEN.String()]; exists {
		metric.Value = fmt.Sprintf("%d", int64(time.Since(mb.start)))
	}
	if metric, exists := mb.metricsMap[type_enums.STATUS.String()]; exists {
		metric.Value = type_enums.RECORD_SUCCESS.String()
	}
	return mb
}

// OnFailure updates the time taken and status metrics for failure
func (mb *MetricBuilder) OnFailure() *MetricBuilder {
	if metric, exists := mb.metricsMap[type_enums.TIME_TAKEN.String()]; exists {
		metric.Value = fmt.Sprintf("%d", int64(time.Since(mb.start)))
	}
	if metric, exists := mb.metricsMap[type_enums.STATUS.String()]; exists {
		metric.Value = type_enums.RECORD_FAILED.String()
	}
	return mb
}

// OnAddMetrics adds additional metrics to the builder, ensuring uniqueness
func (mb *MetricBuilder) OnAddMetrics(metrics ...*protos.Metric) *MetricBuilder {
	for _, newMetric := range metrics {
		mb.metricsMap[newMetric.Name] = newMetric // Ensure uniqueness by overwriting existing ones
	}
	return mb
}

// Build returns the list of unique metrics
func (mb *MetricBuilder) Build() []*protos.Metric {
	uniqueMetrics := make([]*protos.Metric, 0, len(mb.metricsMap))
	for _, metric := range mb.metricsMap {
		uniqueMetrics = append(uniqueMetrics, metric)
	}
	return uniqueMetrics
}
