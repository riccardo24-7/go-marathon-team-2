package metrics

type MetricInfo struct {
	Id      int      // ID в сервисе метрик
	Query   string   // Это название метрики, то самое, которое мы будем хранить в бд
	Folders []string // Список папок, в которых содержится данная метрика
}

type MetricProvider interface {
	GetMetrics() ([]MetricInfo, error)
}

type MetricStorage interface {
	PutMetrics([]MetricInfo) error
}
