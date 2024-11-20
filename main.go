package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type myCollector struct {
  accountBalance prometheus.Gauge
}

func newCollector() *myCollector {
  return &myCollector{
    accountBalance: prometheus.NewGauge(prometheus.GaugeOpts{
      Name: "account_balance",
      Help: "Zadarma account balance metric",
    }),
  }
}

func (c *myCollector) Collect(ch chan<- prometheus.Metric) {
  balance := 5.0
  c.accountBalance.Set(balance)

  ch <- c.accountBalance
}

func (c *myCollector) Describe(ch chan<- *prometheus.Desc) {
  c.accountBalance.Describe(ch)
}
func main() {
  collector := newCollector()
  prometheus.MustRegister(collector)

  http.Handle("/metrics", promhttp.Handler())
  log.Fatal(http.ListenAndServe(":9101", nil))
}
