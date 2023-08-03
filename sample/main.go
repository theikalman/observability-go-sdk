// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package main

import (
	"github.com/AccelByte/observability-go-sdk/metrics"
	"net/http"
	"time"
)

func main() {
	totalSession := metrics.CounterVec(
		"ab_session_total_session",
		"The total number of available session",
		[]string{"game_namespace", "matchpool"},
	)

	metrics.Initialize("test_service")

	go sendCustomPeriodically(totalSession)
	http.Handle("/metrics", metrics.PrometheusHandler())
	http.ListenAndServe(":8080", nil)
}

func sendCustomPeriodically(totalSession metrics.CounterVecMetric) {
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			totalSession.With(map[string]string{"game_namespace": "test", "matchpool": "asdf"}).Add(float64(5))

		case <-quit:
			return
		}
	}
}