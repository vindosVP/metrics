package handlers

import (
	"fmt"
	"net/http"
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

		counterMetrics, err := s.GetAllCounter()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gaugeMetrics, err := s.GetAllGauge()
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

		_, err = w.Write([]byte(html))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}

func counterMetricLines(metrics map[string]int64) []string {
	lines := make([]string, len(metrics))
	index := 0
	for key, val := range metrics {
		line := fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", key, val)
		lines[index] = line
		index++
	}
	return lines
}

func gaugeMetricLines(metrics map[string]float64) []string {
	lines := make([]string, len(metrics))
	index := 0
	for key, val := range metrics {
		line := fmt.Sprintf("<tr><td>%s</td><td>%f</td></tr>", key, val)
		lines[index] = line
		index++
	}
	return lines
}
