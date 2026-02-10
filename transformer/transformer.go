package transformer

type ServerMetric struct {
	Name  string  // Название метрики (например, "memory_usage")
	Value float64 // Значение в байтах
}

var (
	_mb float64 = 1 << 20 // 1024b * 1024b
)

func TransformToMB(mch <-chan ServerMetric) <-chan ServerMetric {
	out := make(chan ServerMetric)

	go func() {
		for m := range mch {
			m.Value = float64(m.Value / _mb)
			out <- m
		}

		close(out)
	}()

	return out
}

type MetricTransformFunc func(in <-chan ServerMetric) <-chan ServerMetric

type MetricsPipeline struct {
	transfn []MetricTransformFunc
}

func (t *MetricsPipeline) Add(tfn MetricTransformFunc) {
	t.transfn = append(t.transfn, tfn)
}

func (t *MetricsPipeline) StartTransform(in <-chan ServerMetric) <-chan ServerMetric {
	if len(t.transfn) == 0 { // или возвращать ошибку вторым значением
		return in
	}

	nextIn := in
	for _, tfn := range t.transfn {
		nextIn = tfn(nextIn)
	}

	return nextIn
}
