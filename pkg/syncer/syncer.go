package syncer

import (
	"fmt"
	"github.com/robfig/cron"
	"go-marathon-team-2/pkg/configuration"
	"go-marathon-team-2/pkg/logger"
	"go-marathon-team-2/pkg/metrics/providers"
	"go-marathon-team-2/pkg/metrics/storages"
	"os"
	"strconv"
	"sync"
)

var config *configuration.Configuration

func StartApp(conf *configuration.Configuration) {
	config = conf

	wgLogger := sync.WaitGroup{}
	go logger.NewLogger(os.Stdout, true, &wgLogger)
	logger.LogChan <- logger.NewMessage(logger.INFO, "Logger init", nil)
	logger.LogChan <- logger.NewMessage(logger.INFO, "Configuration initialized", nil)

	if config.Schedule.Repeat {
		synchronizerBackgroundStart()
	} else {
		synchronizer()
	}
	wgLogger.Wait()
}

func synchronizer() {
	logger.LogChan <- logger.NewMessage(logger.INFO, "Starting new syncer iteration", nil)

	grafanaMetrics, err := providers.NewGrafana(config).GetMetrics()
	if err != nil {
		logger.LogMessageError("An error occurred in the provider Grafana", err)
	}

	apiMetrics, err := providers.NewMetricsApi(config).GetMetrics()
	if err != nil {
		logger.LogMessageError("An error occurred in the provider MetricsApi", err)
	}

	for _, grafMetric := range grafanaMetrics {
		for ind, apiMetric := range apiMetrics {
			if grafMetric.Query == apiMetric.Query {
				apiMetrics[ind].Folders = grafMetric.Folders
			}
		}
	}

	consumer := storages.NewMetricsApi(config)
	err = consumer.PutMetrics(apiMetrics)
	if err != nil {
		logger.LogMessageError("An error occurred in the consumer", err)
	}

	logger.LogChan <- logger.NewMessage(logger.INFO, "Synchronized metrics with success", nil)
	logger.LogChan <- logger.NewMessage(logger.INFO, fmt.Sprintf("Next iteration: days - %d, hours - %d, minutes - %d",
		config.Schedule.Days, config.Schedule.Hours, config.Schedule.Minutes), nil)
}

func synchronizerBackgroundStart() {
	//первый запуск синхронизатора
	logger.LogChan <- logger.NewMessage(logger.INFO, "Synchronizer by schedule starts", nil)
	synchronizer()
	wg := sync.WaitGroup{}
	wg.Add(1)
	schedule := strconv.Itoa(config.Schedule.Days*24+config.Schedule.Hours) +
		"h" + strconv.Itoa(config.Schedule.Minutes) + "m"
	cronJobRunner := cron.New()
	err := cronJobRunner.AddFunc("@every "+schedule, synchronizer)
	if err != nil {
		logger.LogMessageFatal("Schedule is down", err)
	}
	//периодичный запуск синхронизатора
	cronJobRunner.Start()
	wg.Wait()
}
