package statistic

import (
	"time"
)

type LatencyDistribution struct {
	Percentage int
	Latency    time.Duration
}

type Bucket struct {
	Mark      float64
	Count     int
	Frequency float64
}

var (
	pctls = []int{10, 25, 50, 75, 90, 95, 99}
)

func Latencies(latencies []float64) []LatencyDistribution {
	var (
		data = make([]float64, len(pctls))
	)

	j := 0
	for i := 0; i < len(latencies) && j < len(pctls); i++ {
		current := i * 100 / len(latencies)
		if current >= pctls[j] {
			data[j] = latencies[i]
			j++
		}
	}

	res := make([]LatencyDistribution, len(pctls))
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			lat := time.Duration(data[i] * float64(time.Second))
			res[i] = LatencyDistribution{Percentage: pctls[i], Latency: lat}
		}
	}

	return res
}

func Histogram(latencies []float64, slowest, fastest float64) []Bucket {
	var (
		bi      int
		max     int
		bc      = 10
		buckets = make([]float64, bc+1)
		counts  = make([]int, bc+1)
	)

	bs := (slowest - fastest) / float64(bc)
	for i := 0; i < bc; i++ {
		buckets[i] = fastest + bs*float64(i)
	}
	buckets[bc] = slowest
	for i := 0; i < len(latencies); {
		if latencies[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		}
	}

	res := make([]Bucket, len(buckets))
	for i := 0; i < len(buckets); i++ {
		res[i] = Bucket{
			Mark:      buckets[i],
			Count:     counts[i],
			Frequency: float64(counts[i]) / float64(len(latencies)),
		}
	}

	return res
}
