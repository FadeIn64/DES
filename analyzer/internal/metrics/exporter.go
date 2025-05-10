package metrics

import (
	"DAS/models"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	mu *sync.Mutex

	// Основные метрики
	currentLapTime *prometheus.GaugeVec
	avgLapTime     *prometheus.GaugeVec
	segmentPace    *prometheus.GaugeVec
	lapsInSegment  *prometheus.GaugeVec
	trendDirection *prometheus.GaugeVec

	// Дополнительные метрики
	lapDeviation *prometheus.GaugeVec
}

func NewMetricsExporter() *Exporter {
	return &Exporter{
		currentLapTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "f1_current_lap_time_seconds",
				Help: "Current lap time by driver",
			},
			[]string{"driver_number"},
		),
		avgLapTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "f1_avg_lap_time_seconds",
				Help: "Average lap time by driver (excluding pit stops)",
			},
			[]string{"driver_number"},
		),
		segmentPace: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "f1_segment_pace_seconds",
				Help: "Current segment pace by driver",
			},
			[]string{"driver_number", "segment_type"},
		),
		lapsInSegment: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "f1_laps_in_segment_count",
				Help: "Number of laps in current segment",
			},
			[]string{"driver_number"},
		),
		trendDirection: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "f1_lap_trend_direction",
				Help: "Lap trend direction (1=improving, 0=stable, -1=declining)",
			},
			[]string{"driver_number"},
		),
		lapDeviation: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "f1_lap_deviation_percent",
				Help: "Percentage deviation from average lap time",
			},
			[]string{"driver_number"},
		),
		mu: &sync.Mutex{},
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.currentLapTime.Describe(ch)
	e.avgLapTime.Describe(ch)
	e.segmentPace.Describe(ch)
	e.lapsInSegment.Describe(ch)
	e.trendDirection.Describe(ch)
	e.lapDeviation.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.currentLapTime.Collect(ch)
	e.avgLapTime.Collect(ch)
	e.segmentPace.Collect(ch)
	e.lapsInSegment.Collect(ch)
	e.trendDirection.Collect(ch)
	e.lapDeviation.Collect(ch)
}

func (e *Exporter) UpdateMetrics(analysis *models.LapAnalysis) {
	e.mu.Lock()
	defer e.mu.Unlock()

	driverLabel := fmt.Sprintf("%v", analysis.DriverNumber)

	if analysis.CurrentLapTime != 0 {
		e.currentLapTime.WithLabelValues(driverLabel).Set(analysis.CurrentLapTime)
		e.lapDeviation.WithLabelValues(driverLabel).Set(analysis.ComparisonWithAvg)

		// Кодируем тренд в числовое значение
		var trendValue float64
		switch analysis.PositionTrend {
		case "improving":
			trendValue = 1
		case "declining":
			trendValue = -1
		default:
			trendValue = 0
		}
		e.trendDirection.WithLabelValues(driverLabel).Set(trendValue)
	}
	e.avgLapTime.WithLabelValues(driverLabel).Set(analysis.AverageLapTime)
	e.segmentPace.WithLabelValues(driverLabel, "current").Set(analysis.AverageSegmentPace)
	e.lapsInSegment.WithLabelValues(driverLabel).Set(float64(analysis.LapsInSegment))

}
