package collector

type Collector interface {
	Collect() error
}
