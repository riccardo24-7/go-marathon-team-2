package tests

import (
	"go-marathon-team-2/pkg/configuration"
	"go-marathon-team-2/pkg/logger"
	"go-marathon-team-2/pkg/metrics/providers"
	"os"

	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var address = "https://dev.gnivc.ru/tools/metrics-registry/api/v2"

func TestMetricsApiGet_GetIDSuccess(t *testing.T) {
	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	assert.NoError(t, err)
	wg := sync.WaitGroup{}
	go logger.NewLogger(logFile, true, &wg)
	defer logFile.Close()

	config := &configuration.Configuration{}
	config.MetricsAPI.Address = address

	client := providers.NewMetricsApi(config)
	wgTest := sync.WaitGroup{}
	wgTest.Add(1)
	dbMetrics, err := client.GetMetrics()
	assert.NoError(t, err)
	sliceTestInt := []int{31, 2}

	logger.LogChan <- logger.NewMessage(logger.INFO, "Завершение работы теста на чтение метрик из Api", nil)

	assert.Equal(t, sliceTestInt[0], dbMetrics[len(dbMetrics)-1].Id)
	assert.Equal(t, sliceTestInt[1], dbMetrics[len(dbMetrics)-2].Id)
	wg.Wait()
}

func TestMetricsApiGet_GetQuerySuccess(t *testing.T) {
	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	assert.NoError(t, err)
	wg := sync.WaitGroup{}
	go logger.NewLogger(logFile, true, &wg)
	defer logFile.Close()

	config := &configuration.Configuration{}
	config.MetricsAPI.Address = address

	client := providers.NewMetricsApi(config)
	wgTest := sync.WaitGroup{}
	wgTest.Add(1)
	assert.NoError(t, err)
	dbMetrics, err := client.GetMetrics()
	sliceTestInt := "SELECT * FROM metrics_samples;"
	assert.NoError(t, err)
	assert.Equal(t, sliceTestInt, dbMetrics[len(dbMetrics)-1].Query)
}

func TestMetricsApiGet_GetFoldersSuccess(t *testing.T) {
	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	assert.NoError(t, err)
	wg := sync.WaitGroup{}
	go logger.NewLogger(logFile, true, &wg)
	defer logFile.Close()

	config := &configuration.Configuration{}
	config.MetricsAPI.Address = address

	client := providers.NewMetricsApi(config)
	wgTest := sync.WaitGroup{}
	wgTest.Add(1)
	dbMetrics, err := client.GetMetrics()
	assert.NoError(t, err)
	sliceTestInt := []string{"Oracle", "ГИР БО"}
	assert.NoError(t, err)
	assert.Equal(t, sliceTestInt, dbMetrics[len(dbMetrics)-1].Folders)
}
