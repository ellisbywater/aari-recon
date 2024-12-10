package techa

import "math"

type Momentum struct{}

func (m *Momentum) WilliamsR(prices []float64, period int) []float64 {
	// Validate input
	if len(prices) < period {
		return make([]float64, len(prices))
	}

	// Initialize result slice
	williamsR := make([]float64, len(prices))

	// Calculate Williams %R for each valid window
	for i := period - 1; i < len(prices); i++ {
		// Find highest high and lowest low in the lookback period
		highestHigh := math.Inf(-1)
		lowestLow := math.Inf(1)

		for j := i - (period - 1); j <= i; j++ {
			highestHigh = math.Max(highestHigh, prices[j])
			lowestLow = math.Min(lowestLow, prices[j])
		}

		// Prevent division by zero
		if highestHigh == lowestLow {
			williamsR[i] = 0
		} else {
			// Williams %R Formula: %R = (Highest High - Close) / (Highest High - Lowest Low) * -100
			williamsR[i] = ((highestHigh - prices[i]) / (highestHigh - lowestLow)) * -100
		}
	}

	return williamsR
}
