package storages

import (
	"bytes"
	"encoding/json"
	"go-marathon-team-2/pkg/configuration"
	"go-marathon-team-2/pkg/logger"
	"go-marathon-team-2/pkg/metrics"
	"io/ioutil"
	"net/http"
	"strconv"
)

type metricsApi struct {
	config *configuration.Configuration
}

func NewMetricsApi(conf *configuration.Configuration) metrics.MetricStorage {
	return &metricsApi{config: conf}
}

type bodyCtor struct {
	MetricGroups []string `json:"metric_groups"`
}

type respBody struct {
	Id    int    `json:"id"`
	Error string `json:"error"`
}

func (m *metricsApi) PutMetrics(infos []metrics.MetricInfo) error {

	logger.LogMessageInfo("Putting metrics to MetricsApi", nil)
	for _, metric := range infos {
		client := &http.Client{}
		tmpBody := bodyCtor{MetricGroups: metric.Folders}
		someBytes, _ := json.Marshal(tmpBody)

		req, err := http.NewRequest("PUT", m.config.MetricsAPI.Address+"/metrics/"+strconv.Itoa(metric.Id),
			bytes.NewBuffer(someBytes))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		resp, _ := client.Do(req)

		switch resp.StatusCode {
		case http.StatusOK:
			bodyAsBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			var respBody respBody
			err = json.Unmarshal(bodyAsBytes, &respBody)
			if err != nil {
				return err
			}
			if respBody.Id != metric.Id {
				return err
			}
		case http.StatusBadRequest:
			return err
		case http.StatusUnauthorized:
			bodyAsBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			var respBody respBody
			err = json.Unmarshal(bodyAsBytes, &respBody)
			if err != nil {
				return err
			}
			return err
		case http.StatusInternalServerError:
			return err
		}
		resp.Body.Close()
	}
	logger.LogMessageInfo("Success", nil)
	return nil
}
