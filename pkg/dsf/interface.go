package dsf

type DataSource interface {
	GetCpuUsageSample(name DataSourceObjectName) (DataSample, error)
	GetMemoryUsageSample(name DataSourceObjectName) (DataSample, error)
}
