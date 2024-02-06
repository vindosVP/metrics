package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <link rel="stylesheet" type="text/css" href="/assets/main.css">
    <meta charset="UTF-8">
    <title>Live metrics</title>
</head>
<body>
<table>
    <thead>
    <tr>
        <th>Name</th>
        <th>Value</th>
    </tr>
    </thead>
    <tbody>
    %metrics%
    </tbody>
</table>
</body>
</html>`

func List(s MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		counterMetrics, err := s.GetAllCounter(req.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gaugeMetrics, err := s.GetAllGauge(req.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		metricLines := make([]string, 0)
		counterLines := counterMetricLines(counterMetrics)
		gaugeLines := gaugeMetricLines(gaugeMetrics)
		metricLines = append(metricLines, counterLines...)
		metricLines = append(metricLines, gaugeLines...)

		html := strings.Replace(htmlTemplate, "%metrics%", strings.Join(metricLines, ""), -1)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = w.Write([]byte(html))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func counterMetricLines(metrics map[string]int64) []string {
	lines := make([]string, 0, len(metrics))
	keys := make([]string, 0, len(metrics))
	for k := range metrics {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		line := fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", key, metrics[key])
		lines = append(lines, line)
	}
	return lines
}

func gaugeMetricLines(metrics map[string]float64) []string {
	lines := make([]string, 0, len(metrics))
	keys := make([]string, 0, len(metrics))
	for k := range metrics {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		line := fmt.Sprintf("<tr><td>%s</td><td>%.2f</td></tr>", key, metrics[key])
		lines = append(lines, line)
	}
	return lines
}
