package postgresql

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	poolTotalConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "db",
			Subsystem: "pool",
			Name:      "total_connections",
			Help:      "Total number of connections in the pool",
		},
	)

	poolIdleConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "db",
			Subsystem: "pool",
			Name:      "idle_connections",
			Help:      "Number of idle connections in the pool",
		},
	)

	poolAcquiredConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "db",
			Subsystem: "pool",
			Name:      "acquired_connections",
			Help:      "Number of acquired (in-use) connections",
		},
	)

	poolAcquireCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "db",
			Subsystem: "pool",
			Name:      "acquire_total",
			Help:      "Total number of connection acquisitions",
		},
	)

	poolAcquireDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "db",
			Subsystem: "pool",
			Name:      "acquire_duration_seconds",
			Help:      "Time spent acquiring DB connections",
			Buckets:   prometheus.DefBuckets,
		},
	)

	poolCanceledAcquireCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "db",
			Subsystem: "pool",
			Name:      "acquire_canceled_total",
			Help:      "Total number of canceled connection acquire attempts",
		},
	)

	poolEmptyAcquireCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "db",
			Subsystem: "pool",
			Name:      "acquire_empty_total",
			Help:      "Total number of failed acquires due to empty pool",
		},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(
		poolTotalConns,
		poolIdleConns,
		poolAcquiredConns,
		poolAcquireCount,
		poolAcquireDuration,
		poolCanceledAcquireCount,
		poolEmptyAcquireCount,
	)
}

func StartPoolMetrics(pool *pgxpool.Pool, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var prevAcquireCount int64
		var prevCanceledCount int64
		var prevEmptyCount int64

		for range ticker.C {
			stat := pool.Stat()

			poolTotalConns.Set(float64(stat.TotalConns()))
			poolIdleConns.Set(float64(stat.IdleConns()))
			poolAcquiredConns.Set(float64(stat.AcquiredConns()))

			// Counters must be monotonic → use deltas
			acquireCount := stat.AcquireCount()
			poolAcquireCount.Add(float64(acquireCount - prevAcquireCount))
			prevAcquireCount = acquireCount

			poolAcquireDuration.Observe(stat.AcquireDuration().Seconds())

			canceled := stat.CanceledAcquireCount()
			poolCanceledAcquireCount.Add(float64(canceled - prevCanceledCount))
			prevCanceledCount = canceled

			empty := stat.EmptyAcquireCount()
			poolEmptyAcquireCount.Add(float64(empty - prevEmptyCount))
			prevEmptyCount = empty

			// tele.Info(context.Background(), "Pool Metrics",
			// 	"total conns", stat.TotalConns(),
			// 	"idle conns", stat.IdleConns(),
			// 	"acquired conns", acquireCount,
			// 	"Pool Acquire Duration", stat.AcquireDuration().Seconds(),
			// 	"canceled acquire count", canceled,
			// 	"empty acquire count", empty,
			// )
		}
	}()
}
