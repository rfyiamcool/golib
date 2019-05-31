package statistic

import (
	"testing"
)

func TestLatencies(t *testing.T) {
	d := []float64{
		1,
		3,
		3.5,
		20,
		21,
		22,
	}

	t.Log(Latencies(d))
}
