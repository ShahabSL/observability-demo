package main

import (
    "math/rand"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "go.uber.org/zap"
)

var (
    requestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "app_requests_total",
            Help: "Number of requests to each endpoint",
        },
        []string{"endpoint", "status"},
    )

    activeUsers = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "app_active_users",
            Help: "Simulated number of active users",
        },
    )

    logger *zap.Logger
)

func main() {
    // Initialize logger
    logger, _ = zap.NewProduction()
    defer logger.Sync()

    // Register Prometheus metrics
    prometheus.MustRegister(requestCounter)
    prometheus.MustRegister(activeUsers)

    // Simulate some "active users" metric
    go func() {
        for {
            activeUsers.Set(float64(rand.Intn(100)))
            time.Sleep(5 * time.Second)
        }
    }()

    // Basic homepage
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        requestCounter.WithLabelValues("home", "success").Inc()
        logger.Info("Homepage accessed")
        w.Write([]byte("Hello! This is our observable app"))
    })

    // Endpoint that randomly succeeds/fails
    http.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
        if rand.Float64() < 0.5 {
            requestCounter.WithLabelValues("random", "success").Inc()
            logger.Info("Random endpoint success")
            w.Write([]byte("Success!"))
        } else {
            requestCounter.WithLabelValues("random", "error").Inc()
            logger.Error("Random endpoint failed")
            http.Error(w, "Random failure!", http.StatusInternalServerError)
        }
    })

    // Prometheus metrics endpoint
    http.Handle("/metrics", promhttp.Handler())

    logger.Info("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}