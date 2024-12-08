package techa

func lowestLow(values []float64) float64 {
	lowest := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < lowest {
			lowest = values[i]
		}
	}
	return lowest
}

func highestHigh(values []float64) float64 {
	highest := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > highest {
			highest = values[i]
		}
	}
	return highest
}

func highestAndLowest(values []float64) (float64, float64) {
	highest := values[0]
	lowest := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > highest {
			highest = values[i]
		} else if values[i] < lowest {
			lowest = values[i]
		}
	}
	return highest, lowest
}
