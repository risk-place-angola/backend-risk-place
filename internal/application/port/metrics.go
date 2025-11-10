package port

type Metrics interface {
	Increment(metric string, value int)
	Observe(metric string, value float64)
}
