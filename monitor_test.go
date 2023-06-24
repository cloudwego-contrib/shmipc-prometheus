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
	"testing"

	"github.com/cloudwego/shmipc-go"
	"github.com/stretchr/testify/assert"
)

func TestNewPrometheusMonitor(t *testing.T) {
	addr := "localhost:9090"
	path := "/metrics"
	monitor := NewPrometheusMonitor(addr, path)

	assert.NotNil(t, monitor.receiveSyncEventCount)
	assert.NotNil(t, monitor.sendSyncEventCount)
	assert.NotNil(t, monitor.outFlowBytes)
	assert.NotNil(t, monitor.inFlowBytes)
	assert.NotNil(t, monitor.sendQueueCount)
	assert.NotNil(t, monitor.receiveQueueCount)
	assert.NotNil(t, monitor.allocShmErrorCount)
	assert.NotNil(t, monitor.fallbackWriteCount)
	assert.NotNil(t, monitor.fallbackReadCount)
	assert.NotNil(t, monitor.eventConnErrorCount)
	assert.NotNil(t, monitor.queueFullErrorCount)
	assert.NotNil(t, monitor.activeStreamCount)
	assert.NotNil(t, monitor.hotRestartSuccessCount)
	assert.NotNil(t, monitor.hotRestartErrorCount)
	assert.NotNil(t, monitor.capacityOfShareMemory)
	assert.NotNil(t, monitor.allInUsedShareMemory)
	assert.NotNil(t, monitor.MonitorInfo)
}

func TestOnEmitSessionMetricsAndFlush(t *testing.T) {
	addr := "localhost:9090"
	path := "/metrics"
	monitor := NewPrometheusMonitor(addr, path)

	performanceMetrics := shmipc.PerformanceMetrics{
		ReceiveSyncEventCount: 10,
		SendSyncEventCount:    20,
		OutFlowBytes:          30,
		InFlowBytes:           40,
		SendQueueCount:        50,
		ReceiveQueueCount:     60,
	}

	stabilityMetrics := shmipc.StabilityMetrics{
		AllocShmErrorCount:     1,
		FallbackWriteCount:     2,
		FallbackReadCount:      3,
		EventConnErrorCount:    4,
		QueueFullErrorCount:    5,
		ActiveStreamCount:      6,
		HotRestartSuccessCount: 7,
		HotRestartErrorCount:   8,
	}

	shareMemoryMetrics := shmipc.ShareMemoryMetrics{
		CapacityOfShareMemoryInBytes: 1024 * 1024,
		AllInUsedShareMemoryInBytes:  512 * 1024,
	}

	monitor.OnEmitSessionMetrics(performanceMetrics, stabilityMetrics, shareMemoryMetrics, nil)

	assert.Equal(t, float64(10), monitor.MonitorInfo["receiveSyncEventCount"])
	assert.Equal(t, float64(20), monitor.MonitorInfo["sendSyncEventCount"])
	assert.Equal(t, float64(30), monitor.MonitorInfo["outFlowBytes"])
	assert.Equal(t, float64(40), monitor.MonitorInfo["inFlowBytes"])
	assert.Equal(t, float64(50), monitor.MonitorInfo["sendQueueCount"])
	assert.Equal(t, float64(60), monitor.MonitorInfo["receiveQueueCount"])
	assert.Equal(t, float64(1), monitor.MonitorInfo["allocShmErrorCount"])
	assert.Equal(t, float64(2), monitor.MonitorInfo["fallbackWriteCount"])
	assert.Equal(t, float64(3), monitor.MonitorInfo["fallbackReadCount"])
	assert.Equal(t, float64(4), monitor.MonitorInfo["eventConnErrorCount"])
	assert.Equal(t, float64(5), monitor.MonitorInfo["queueFullErrorCount"])
	assert.Equal(t, float64(6), monitor.MonitorInfo["activeStreamCount"])
	assert.Equal(t, float64(7), monitor.MonitorInfo["hotRestartSuccessCount"])
	assert.Equal(t, float64(8), monitor.MonitorInfo["hotRestartErrorCount"])
	assert.Equal(t, float64(1024*1024), monitor.MonitorInfo["capacityOfShareMemory"])
	assert.Equal(t, float64(512*1024), monitor.MonitorInfo["allInUsedShareMemory"])

	// flush the metrics to the Prometheus server
	err := monitor.Flush()
	assert.Nil(t, err)
}
