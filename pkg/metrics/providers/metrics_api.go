package providers

import (
	"encoding/json"
	"go-marathon-team-2/pkg/configuration"
	"go-marathon-team-2/pkg/logger"
	"go-marathon-team-2/pkg/metrics"
	"io/ioutil"
	"net/http"
)

type metricsApi struct {
	config *configuration.Configuration
}

func NewMetricsApi(conf *configuration.Configuration) metrics.MetricProvider {
	return &metricsApi{config: conf}
}

type webJSON struct {
	ID           int      `json:"id"`
	Query        string   `json:"query"`
	MetricGroups []string `json:"metric_groups"`
}

func (m *metricsApi) GetMetrics() ([]metrics.MetricInfo, error) {

	logger.LogMessageInfo("Getting metrics from MetricsApi", nil)
	var result []metrics.MetricInfo
	response, err := http.Get(m.config.MetricsAPI.Address + "/metrics")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var webresp []webJSON
	err = json.Unmarshal(bodyBytes, &webresp)
	if err != nil {
		return nil, err
	}

	for _, elem := range webresp {
		result = append(result, metrics.MetricInfo{
			Id:      elem.ID,
			Query:   elem.Query,
			Folders: elem.MetricGroups,
		})
	}

	logger.LogMessageInfo("Success", nil)
	return result, nil
}
