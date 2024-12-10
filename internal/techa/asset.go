package techa

import "time"

type Asset struct {
	Name    string
	Date    []time.Time
	Opening []float64
	Closing []float64
	High    []float64
	Low     []float64
	Volume  []float64
}
