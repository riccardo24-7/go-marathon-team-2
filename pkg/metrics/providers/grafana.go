package providers

import (
	"encoding/json"
	"go-marathon-team-2/pkg/configuration"
	"go-marathon-team-2/pkg/logger"
	"go-marathon-team-2/pkg/metrics"
	"io/ioutil"
	"net/http"
	"sync"
)

type grafana struct {
	config *configuration.Configuration
}

func NewGrafana(conf *configuration.Configuration) metrics.MetricProvider {
	return &grafana{config: conf}
}

type grafanaStruct struct {
	Uid       string `json:"uid"`
	SomeType  string `json:"type"`
	FolderUid string `json:"folderUid"`
	Title     string `json:"title"`
}

type dashboardStruct struct {
	Db struct {
		Panels []struct {
			Targets []struct {
				RawSql string `json:"rawSql"`
				Expr   string `json:"expr"`
			} `json:"targets"`
		} `json:"panels"`
	} `json:"dashboard"`
}

type grafanaDashboard struct {
	Uid        string
	FolderName string
}

func (g *grafana) GetMetrics() ([]metrics.MetricInfo, error) {

	logger.LogMessageInfo("Getting metrics from Grafana", nil)
	defer logger.LogMessageInfo("Success", nil)

	client := &http.Client{}
	// To return
	var result []metrics.MetricInfo
	// Chans
	chMetricsInternal := make(chan [2]string)
	chMetricsError := make(chan error)
	// Maps
	metricsMap := make(map[string][]string)
	metricsMapReversed := make(map[string][]string)
	// Waitgroups
	wgMetrics := sync.WaitGroup{}
	// Делаем запрос, чтобы получить дашборды и папки
	req, err := http.NewRequest("GET", g.config.Grafana.Address+"/api/search", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+g.config.Grafana.Key)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteAnswer, _ := ioutil.ReadAll(resp.Body)

	var grafanaMap []grafanaStruct
	err = json.Unmarshal(byteAnswer, &grafanaMap)
	if err != nil {
		return nil, err
	}

	mappedGrafana := responseFormatter(grafanaMap) // Форматируем и добавляем название папки в FolderUid, при этом отсеиваем объекты папок
	// Асинхронно считываем метрики из графаны по Uid
	go func() {
		for {
			select {
			case metric := <-chMetricsInternal:
				if metric[0] != "" {
					metricsMap[metric[0]] = append(metricsMap[metric[0]], metric[1])
					wgMetrics.Done()
				}
			}
		}
	}()

	for _, db := range mappedGrafana {
		if db.Uid != "" {
			wgMetrics.Add(1)
			go asyncGetMetricsByUid(
				&wgMetrics,
				g.config.Grafana.Address,
				g.config.Grafana.Key,
				db,
				chMetricsInternal,
				chMetricsError,
			)
		}
	}
	wgMetrics.Wait()

	// Проверяем пришла ли ошибка
	if len(chMetricsError) != 0 {
		err = <-chMetricsError
	}
	if err != nil {
		return nil, err
	}

	// Reverting - меняем в мапе местами названия папок и названия запросов
	flag := true
	for key, values := range metricsMap {
		for _, value := range values {
			for _, elementInMap := range metricsMapReversed[value] {
				if elementInMap == key {
					flag = false
				}
			}
			if flag {
				metricsMapReversed[value] = append(metricsMapReversed[value], key)
			}
			flag = true
		}
	}

	// Creating metrics - заносим все в конструктор
	for key, values := range metricsMapReversed {
		result = append(result, metrics.MetricInfo{
			Id:      0,
			Query:   key,
			Folders: values,
		})
	}

	// Закрываем каналы
	close(chMetricsInternal)
	close(chMetricsError)

	return result, nil
}

func asyncGetMetricsByUid(wg *sync.WaitGroup, address, key string, db grafanaDashboard, ch chan [2]string, chError chan error) {
	defer wg.Done()

	client := &http.Client{}

	req, err := http.NewRequest("GET", address+"/api/dashboards/uid/"+db.Uid, nil)
	if err != nil {
		chError <- err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	response, _ := client.Do(req)
	defer response.Body.Close()

	byteAnswer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		chError <- err
	}
	var dashboardsMap dashboardStruct
	err = json.Unmarshal(byteAnswer, &dashboardsMap)
	if err != nil {
		chError <- err
	}

	for _, panel := range dashboardsMap.Db.Panels {
		wg.Add(1)
		if panel.Targets[0].RawSql == "" {
			ch <- [2]string{db.FolderName, panel.Targets[0].Expr}
		} else {
			ch <- [2]string{db.FolderName, panel.Targets[0].RawSql}
		}
	}
}

func newGrafanaBoard(uid, folderName string) grafanaDashboard {
	return grafanaDashboard{
		Uid:        uid,
		FolderName: folderName,
	}
}

func responseFormatter(g []grafanaStruct) []grafanaDashboard {

	rez := make([]grafanaDashboard, 10)

	for _, item := range g {
		if item.SomeType == "dash-db" {
			if item.FolderUid == "" {
				rez = append(rez, newGrafanaBoard(item.Uid, "general"))
			} else {
				for _, folder := range g {
					if item.FolderUid == folder.Uid {
						rez = append(rez, newGrafanaBoard(item.Uid, folder.Title))
					}
				}
			}
		}
	}
	return rez
}
