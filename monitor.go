/*
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package shmipc

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/shmipc-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMonitor struct {
	receiveSyncEventCount prometheus.Gauge
	sendSyncEventCount    prometheus.Gauge
	outFlowBytes          prometheus.Gauge
	inFlowBytes           prometheus.Gauge
	sendQueueCount        prometheus.Gauge
	receiveQueueCount     prometheus.Gauge

	allocShmErrorCount     prometheus.Gauge
	fallbackWriteCount     prometheus.Gauge
	fallbackReadCount      prometheus.Gauge
	eventConnErrorCount    prometheus.Gauge
	queueFullErrorCount    prometheus.Gauge
	activeStreamCount      prometheus.Gauge
	hotRestartSuccessCount prometheus.Gauge
	hotRestartErrorCount   prometheus.Gauge

	capacityOfShareMemory prometheus.Gauge
	allInUsedShareMemory  prometheus.Gauge

	MonitorInfo map[string]float64
}

func NewPrometheusMonitor(addr, path string) *PrometheusMonitor {
	registry := prometheus.NewRegistry()

	http.Handle(path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal("Unable to start a promhttp server, err: " + err.Error())
		}
	}()

	receiveSyncEventCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "receive_sync_event_count",
		Help: "The SyncEvent count that session had received",
	})
	sendSyncEventCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "send_sync_event_count",
		Help: "The SyncEvent count that session had sent",
	})
	outFlowBytes := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "out_flow_bytes",
		Help: "The out flow in bytes that session had sent",
	})
	inFlowBytes := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "in_flow_bytes",
		Help: "The in flow in bytes that session had receive",
	})
	sendQueueCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "send_queue_count",
		Help: "The pending count of send queue",
	})
	receiveQueueCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "receive_queue_count",
		Help: "The pending count of receive queue",
	})
	allocShmErrorCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "alloc_shm_error_count",
		Help: "The error count of allocating share memory",
	})
	fallbackWriteCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fallback_write_count",
		Help: "The count of the fallback data write to unix/tcp connection",
	})
	fallbackReadCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fallback_read_count",
		Help: "The error count of receiving fallback data from unix/tcp connection every period",
	})
	eventConnErrorCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "event_conn_error_count",
		Help: "The error count of unix/tcp connection which usually happened in that the peer's process exit(crashed or other reason)",
	})
	queueFullErrorCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "queue_full_error_count",
		Help: "The error count due to the IO-Queue(SendQueue or ReceiveQueue) is full which usually happened in that the peer was busy",
	})
	activeStreamCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_stream_count",
		Help: "Current all active stream count",
	})
	hotRestartSuccessCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hot_restart_success_count",
		Help: "The successful count of hot restart",
	})
	hotRestartErrorCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hot_restart_error_count",
		Help: "The failed count of hot restart",
	})
	capacityOfShareMemory := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "capacity_of_share_memory",
		Help: "The capacity of the share memory in bytes",
	})
	allInUsedShareMemory := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "all_in_used_share_memory",
		Help: "The amount of share memory in bytes that is currently in use",
	})

	registry.MustRegister(
		receiveSyncEventCount,
		sendSyncEventCount,
		outFlowBytes,
		inFlowBytes,
		sendQueueCount,
		receiveQueueCount,
		allocShmErrorCount,
		fallbackWriteCount,
		fallbackReadCount,
		eventConnErrorCount,
		queueFullErrorCount,
		activeStreamCount,
		hotRestartSuccessCount,
		hotRestartErrorCount,
		capacityOfShareMemory,
		allInUsedShareMemory,
	)

	return &PrometheusMonitor{
		receiveSyncEventCount:  receiveSyncEventCount,
		sendSyncEventCount:     sendSyncEventCount,
		outFlowBytes:           outFlowBytes,
		inFlowBytes:            inFlowBytes,
		sendQueueCount:         sendQueueCount,
		receiveQueueCount:      receiveQueueCount,
		allocShmErrorCount:     allocShmErrorCount,
		fallbackWriteCount:     fallbackWriteCount,
		fallbackReadCount:      fallbackReadCount,
		eventConnErrorCount:    eventConnErrorCount,
		queueFullErrorCount:    queueFullErrorCount,
		activeStreamCount:      activeStreamCount,
		hotRestartSuccessCount: hotRestartSuccessCount,
		hotRestartErrorCount:   hotRestartErrorCount,
		capacityOfShareMemory:  capacityOfShareMemory,
		allInUsedShareMemory:   allInUsedShareMemory,
		MonitorInfo:            make(map[string]float64),
	}
}

// OnEmitSessionMetrics was called by shmipc-go Session with periodically.
func (p *PrometheusMonitor) OnEmitSessionMetrics(performanceMetrics shmipc.PerformanceMetrics, stabilityMetrics shmipc.StabilityMetrics, shareMemoryMetrics shmipc.ShareMemoryMetrics, session *shmipc.Session) {
	p.receiveSyncEventCount.Set(float64(performanceMetrics.ReceiveSyncEventCount))
	p.sendSyncEventCount.Set(float64(performanceMetrics.SendSyncEventCount))
	p.outFlowBytes.Set(float64(performanceMetrics.OutFlowBytes))
	p.inFlowBytes.Set(float64(performanceMetrics.InFlowBytes))
	p.sendQueueCount.Set(float64(performanceMetrics.SendQueueCount))
	p.receiveQueueCount.Set(float64(performanceMetrics.ReceiveQueueCount))

	p.allocShmErrorCount.Set(float64(stabilityMetrics.AllocShmErrorCount))
	p.fallbackWriteCount.Set(float64(stabilityMetrics.FallbackWriteCount))
	p.fallbackReadCount.Set(float64(stabilityMetrics.FallbackReadCount))
	p.eventConnErrorCount.Set(float64(stabilityMetrics.EventConnErrorCount))
	p.queueFullErrorCount.Set(float64(stabilityMetrics.QueueFullErrorCount))
	p.activeStreamCount.Set(float64(stabilityMetrics.ActiveStreamCount))
	p.hotRestartSuccessCount.Set(float64(stabilityMetrics.HotRestartSuccessCount))
	p.hotRestartErrorCount.Set(float64(stabilityMetrics.HotRestartErrorCount))

	p.capacityOfShareMemory.Set(float64(shareMemoryMetrics.CapacityOfShareMemoryInBytes))
	p.allInUsedShareMemory.Set(float64(shareMemoryMetrics.AllInUsedShareMemoryInBytes))

	p.MonitorInfo["receiveSyncEventCount"] = float64(performanceMetrics.ReceiveSyncEventCount)
	p.MonitorInfo["sendSyncEventCount"] = float64(performanceMetrics.SendSyncEventCount)
	p.MonitorInfo["outFlowBytes"] = float64(performanceMetrics.OutFlowBytes)
	p.MonitorInfo["inFlowBytes"] = float64(performanceMetrics.InFlowBytes)
	p.MonitorInfo["sendQueueCount"] = float64(performanceMetrics.SendQueueCount)
	p.MonitorInfo["receiveQueueCount"] = float64(performanceMetrics.ReceiveQueueCount)

	p.MonitorInfo["allocShmErrorCount"] = float64(stabilityMetrics.AllocShmErrorCount)
	p.MonitorInfo["fallbackWriteCount"] = float64(stabilityMetrics.FallbackWriteCount)
	p.MonitorInfo["fallbackReadCount"] = float64(stabilityMetrics.FallbackReadCount)
	p.MonitorInfo["eventConnErrorCount"] = float64(stabilityMetrics.EventConnErrorCount)
	p.MonitorInfo["queueFullErrorCount"] = float64(stabilityMetrics.QueueFullErrorCount)
	p.MonitorInfo["activeStreamCount"] = float64(stabilityMetrics.ActiveStreamCount)
	p.MonitorInfo["hotRestartSuccessCount"] = float64(stabilityMetrics.HotRestartSuccessCount)
	p.MonitorInfo["hotRestartErrorCount"] = float64(stabilityMetrics.HotRestartErrorCount)

	p.MonitorInfo["capacityOfShareMemory"] = float64(shareMemoryMetrics.CapacityOfShareMemoryInBytes)
	p.MonitorInfo["allInUsedShareMemory"] = float64(shareMemoryMetrics.AllInUsedShareMemoryInBytes)
}

// Flush metrics to log file
func (p *PrometheusMonitor) Flush() error {
	f, err := os.Create(fmt.Sprintf("MonitorInfo_%s.log", time.Now().Format("20060102150405")))
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}
	defer f.Close()

	for key, value := range p.MonitorInfo {
		fmt.Fprintf(f, "%s: %f\n", key, value)
	}

	return nil
}
