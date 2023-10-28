package metrics

import (
	"contrib.go.opencensus.io/exporter/prometheus"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"time"
)

var (
	TagKeyMethod, _ = tag.NewKey("method")
	TagKeyPath, _   = tag.NewKey("path")
)

func SetupOpenCensus() *prometheus.Exporter {

	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "demo",
	})
	if err != nil {
		log.Fatalf("Failed to create the Prometheus stats exporter: %v", err)
	}

	view.RegisterExporter(exporter)
	view.SetReportingPeriod(1 * time.Second)

	latencyView := &view.View{
		Name: "http/latency",
		//Measure:     stats.HTTPServerRequestBytes,
		Description: "Latency distribution of HTTP requests",
		TagKeys:     []tag.Key{TagKeyMethod, TagKeyPath},
		Aggregation: view.Distribution(0, 100, 300, 900, 1800, 2700, 3600, 4500, 5400, 6300, 7200, 8100, 9000),
	}

	_ = view.Register(latencyView)

	return exporter
}
