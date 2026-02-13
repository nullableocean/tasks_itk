package transformer

import (
	"fmt"
	"testing"
)

func logRawMetric(in <-chan *ServerMetric) <-chan *ServerMetric {
	o := make(chan *ServerMetric)

	go func() {
		defer close(o)

		for m := range in {
			fmt.Printf("RAW: name: %s, value: %v\n", m.Name, m.Value)
			o <- m
		}
	}()

	return o
}

func Test_Transformer(t *testing.T) {
	rawMetricsCh := make(chan *ServerMetric)

	go func() {
		defer close(rawMetricsCh)

		for i := range 10 {
			b := 1 << (20 + i + 1)
			testM := &ServerMetric{
				Name:  fmt.Sprintf("node_%d_memusage", i),
				Value: float64(b),
			}

			rawMetricsCh <- testM
		}
	}()
	metricTrans := &MetricsPipeline{}

	metricTrans.Add(logRawMetric)
	metricTrans.Add(TransformToMB)

	apiCh := metricTrans.StartTransform(rawMetricsCh)

	fmt.Println()
	fmt.Println()

	for metric := range apiCh {
		fmt.Println(metric.Name, "MB: ", metric.Value)
	}
}
