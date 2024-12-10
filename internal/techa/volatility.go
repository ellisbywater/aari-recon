package techa

import (
	"fmt"
	"math"
)

type Volatility struct{}

// CalculateRSI computes the Relative Strength Index for a given slice of prices
// and a specified period (typically 14 days)
func (v *Volatility) RSI(prices []float64, period int) []float64 {
	// Validate input
	if len(prices) < period {
		return nil
	}

	// Slice to store RSI values
	rsiValues := make([]float64, len(prices))

	// Calculate initial average gains and losses over the first period
	var avgGain, avgLoss float64
	for i := 1; i < period; i++ {
		change := prices[i] - prices[i-1]
		if change >= 0 {
			avgGain += change
		} else {
			avgLoss -= change
		}
	}

	// Initial averages
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Calculate relative strength
	var rs float64
	if avgLoss != 0 {
		rs = avgGain / avgLoss
	}

	// First RSI calculation
	rsiValues[period-1] = 100.0 - (100.0 / (1.0 + rs))

	// Subsequent RSI calculations using Wilder's smoothing method
	for i := period; i < len(prices); i++ {
		change := prices[i] - prices[i-1]

		// Separate gains and losses
		gain := math.Max(change, 0)
		loss := math.Abs(math.Min(change, 0))

		// Smoothed moving averages
		avgGain = ((avgGain * float64(period-1)) + gain) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + loss) / float64(period)

		// Relative Strength
		if avgLoss != 0 {
			rs = avgGain / avgLoss
		} else {
			rs = 0
		}

		// RSI calculation
		rsiValues[i] = 100.0 - (100.0 / (1.0 + rs))
	}

	return rsiValues
}

func (v *Volatility) BollingerBands(prices []float64, period int, multiplier float64) ([]float64, []float64, []float64) {
	smaVals := make([]float64, len(prices))
	stdVals := make([]float64, len(prices))
	upperVals := make([]float64, len(prices))
	lowerVals := make([]float64, len(prices))

	for i := 0; i < len(prices); i++ {
		if i >= period-1 {
			smaVals[i] = prices[i]
			stdVals[i] = prices[i]
			upperVals[i] = prices[i]
			lowerVals[i] = prices[i]
		} else {
			smaVals[i] = calculateSMASnapshot(prices[i-period : i])
			stdVals[i] = calculateStdDev(prices[i-period:i], period, smaVals[i])
			upperVals[i] = calculateUpperBand(smaVals[i], stdVals[i], multiplier)
			lowerVals[i] = calculateLowerBand(smaVals[i], stdVals[i], multiplier)
		}
	}
	return smaVals, upperVals, lowerVals
}

// calculateStochasticOscillator computes the Stochastic Oscillator
func calculateStochasticOscillator(values []float64, period int) (float64, error) {
	// Validate input
	if len(values) < period {
		return 0, fmt.Errorf("insufficient RSI values")
	}

	// Find the current RSI value (last in the slice)
	currentRSI := values[len(values)-1]

	// Find lowest and highest RSI in the period
	lowestRSI := currentRSI
	highestRSI := currentRSI
	for i := len(values) - period; i < len(values); i++ {
		lowestRSI = math.Min(lowestRSI, values[i])
		highestRSI = math.Max(highestRSI, values[i])
	}

	// Calculate StochRSI
	if highestRSI == lowestRSI {
		return 100, nil // Avoid division by zero
	}

	stochRSI := ((currentRSI - lowestRSI) / (highestRSI - lowestRSI)) * 100

	return stochRSI, nil
}

// StochRSI calculates the Stochastic RSI for a given slice of prices
// Parameters:
// - prices: Slice of float64 representing price data
// - rsiPeriod: Number of periods for RSI calculation
// - stochPeriod: Number of periods for Stochastic calculation
func (v *Volatility) StochRSI(prices []float64, rsiPeriod, stochPeriod int) (float64, error) {
	// Validate input
	if len(prices) < rsiPeriod+stochPeriod {
		return 0, fmt.Errorf("insufficient data points")
	}

	// Calculate RSI values
	rsiValues := v.RSI(prices, rsiPeriod)

	// Calculate StochRSI
	return calculateStochasticOscillator(rsiValues, stochPeriod)
}
