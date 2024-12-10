package techa

import (
	"fmt"
	"math"
)

type Trends struct{}

func calculateSMASnapshot(prices []float64) float64 {
	sum := 0.0
	for i := 0; i < len(prices); i++ {
		sum += prices[i]
	}
	return sum / float64(len(prices))
}

func (trends *Trends) SMA(prices []float64, period int) ([]float64, error) {
	// Validate input
	if period <= 0 {
		return nil, fmt.Errorf("window size must be a positive integer")
	}
	if len(prices) < period {
		return nil, fmt.Errorf("data length must be at least the window size")
	}
	// Initialize result slice
	sma := make([]float64, len(prices)-period+1)

	// Calculate SMA for each window
	for i := 0; i <= len(prices)-period; i++ {
		// Take the slice of current window
		window := prices[i : i+period]

		// Calculate sum of the window
		sum := 0.0
		for _, value := range window {
			sum += value
		}

		// Calculate average
		sma[i] = sum / float64(period)
	}

	return sma, nil
}

func (trends *Trends) TrueRange(high, low, close []float64) []float64 {
	trValues := make([]float64, len(close))

	for i := 1; i < len(close); i++ {
		// True Range is the maximum of:
		// 1. High - Low
		// 2. Abs(High - Previous Close)
		// 3. Abs(Low - Previous Close)
		tr1 := high[i] - low[i]
		tr2 := math.Abs(high[i] - close[i-1])
		tr3 := math.Abs(low[i] - close[i-1])

		trValues[i] = math.Max(math.Max(tr1, tr2), tr3)
	}

	return trValues
}

func (trends *Trends) AvgTrueRange(trValues []float64, period int) []float64 {
	atrValues := make([]float64, len(trValues))

	sum := 0.0
	for i := 1; i <= period; i++ {
		sum += trValues[i]
	}
	atrValues[period] = sum / float64(period)

	// Subsequent ATR values use smoothing
	for i := period + 1; i < len(trValues); i++ {
		atrValues[i] = ((atrValues[i-1] * float64(period-1)) + trValues[i]) / float64(period)
	}

	return atrValues
}

func (trends *Trends) EMA(data []float64, period int) []float64 {
	if len(data) == 0 || period <= 0 {
		return []float64{}
	}

	// Ensure period doesn't exceed data length
	if period > len(data) {
		period = len(data)
	}

	// Calculate smoothing factor
	smoothing := 2.0 / float64(period+1)

	// Initialize the result slice
	ema := make([]float64, len(data))

	// Calculate initial SMA for the first period
	var initialSMA float64
	for i := 0; i < period; i++ {
		initialSMA += data[i]
	}
	initialSMA /= float64(period)

	// First EMA is the initial SMA
	ema[period-1] = initialSMA

	// Calculate subsequent EMAs
	for i := period; i < len(data); i++ {
		ema[i] = (data[i]-ema[i-1])*smoothing + ema[i-1]
	}

	return ema
}

func (trends *Trends) DEMA(prices []float64, period int) []float64 {
	ema1 := trends.EMA(prices, period)
	ema2 := trends.EMA(ema1, period)

	dema := make([]float64, len(prices))

	for i := 2*period - 2; i < len(prices); i++ {
		dema[i] = 2*ema1[i] - ema2[i]
	}
	return dema
}

func (trends *Trends) TREMA(prices []float64, period int) []float64 {
	smoothing := 2.0 / float64(period+1)
	tema := make([]float64, len(prices))

	ema1 := make([]float64, len(prices))
	ema1[0] = prices[0]
	for i := 1; i < len(prices); i++ {
		ema1[i] = prices[i]*smoothing + ema1[i-1]*(1-smoothing)
	}

	dema := trends.DEMA(prices, period)

	// calculate triple ema
	ema3 := make([]float64, len(prices))
	ema3[0] = dema[0]
	for i := 1; i < len(prices); i++ {
		ema3[i] = dema[i]*smoothing + ema3[i-1]*(1-smoothing)
	}

	// Calculate TEMA
	for i := period - 1; i < len(prices); i++ {
		tema[i] = 3*ema1[i] - 3*dema[i] + ema3[i]
	}

	return tema
}

func (trends *Trends) MACD(prices []float64, fast int, slow int, signal int) ([]float64, []float64, []float64) {
	shortEMA := make([]float64, len(prices))
	longEMA := make([]float64, len(prices))
	macdvalue := make([]float64, len(prices))
	signals := make([]float64, len(prices))
	delta := make([]float64, len(prices))

	shortEMA[0] = prices[0]
	longEMA[0] = prices[0]

	for i := 1; i < len(prices); i++ {
		shortEMA[i] = (prices[i]*2 + shortEMA[i-1]*(float64(fast)-2)) / float64(fast)
		longEMA[i] = (prices[i]*2 + longEMA[i-1]*(float64(slow)-2)) / float64(slow)

		// Calculate MACD line
		macdvalue[i] = longEMA[i] - shortEMA[i]

		// Calculate Signal line
		signals[i] = macdvalue[i]*2 + signals[i-1]*(float64(signal)-2)/float64(signal)

		// Calculate Convergence/Divergence
		delta[i] = macdvalue[i] - signals[i]
	}

	return macdvalue, signals, delta
}

type SuperTrendResult struct {
	SuperTrend float64
	Trend      int
}

func (trends *Trends) SuperTrend(high, low, close []float64, period int, multiplier float64) []SuperTrendResult {
	if len(high) != len(low) || len(high) != len(close) {
		return nil
	}

	// Calculate True Range
	trValues := trends.TrueRange(high, low, close)

	// Calculate Average True Range (ATR)
	atrValues := trends.AvgTrueRange(trValues, period)

	// Initialize result slice
	results := make([]SuperTrendResult, len(close))

	// Temporary variables for SuperTrend calculation
	var upperBand, lowerBand float64
	var prevSuperTrend float64
	var trend int

	for i := period; i < len(close); i++ {
		// Calculate basic bands
		basicUpperBand := ((high[i] + low[i]) / 2) + (multiplier * atrValues[i])
		basicLowerBand := ((high[i] + low[i]) / 2) - (multiplier * atrValues[i])

		// SuperTrend calculation
		if basicUpperBand < upperBand || close[i-1] > upperBand {
			upperBand = basicUpperBand
		}

		if basicLowerBand > lowerBand || close[i-1] < lowerBand {
			lowerBand = basicLowerBand
		}

		// Determine trend
		if close[i] <= upperBand {
			// Potential uptrend
			if prevSuperTrend != upperBand {
				prevSuperTrend = upperBand
				trend = 1
			}
		} else {
			// Potential downtrend
			if prevSuperTrend != lowerBand {
				prevSuperTrend = lowerBand
				trend = -1
			}
		}

		// Store results
		results[i] = SuperTrendResult{
			SuperTrend: prevSuperTrend,
			Trend:      trend,
		}
	}
	return results
}

// CalculateTRIX computes the TRIX indicator for a given slice of prices
// Parameters:
// - prices: Input price series (typically closing prices)
// - period: The smoothing period for EMA calculations
// Returns:
// - A slice of TRIX values corresponding to the input prices
func (trends *Trends) TRIX(prices []float64, period int) []float64 {
	if len(prices) < period {
		return []float64{}
	}
	// First EMA (Single smoothing)
	ema1 := trends.EMA(prices, period)

	// Second EMA (Double smoothing)
	ema2 := trends.EMA(ema1, period)

	// Third EMA (Triple smoothing)
	ema3 := trends.EMA(ema2, period)

	// Calculate TRIX: Percentage change of the triple smoothed EMA
	trix := make([]float64, len(ema3))
	for i := 1; i < len(ema3); i++ {
		if ema3[i-1] != 0 {
			trix[i] = (ema3[i] - ema3[i-1]) / ema3[i-1] * 100
		}
	}
	return trix
}

// AroonIndicator calculates the Aroon Up and Down indicators for a given slice of prices
// period: the number of periods to analyze
// prices: a slice of float64 representing price data (typically closing prices)
// Returns: two slices of float64 - Aroon Up and Aroon Down values
func (trends *Trends) Aroon(period int, prices []float64) ([]float64, []float64) {
	if len(prices) < period {
		return nil, nil
	}

	// Prepare output slices
	aroonUp := make([]float64, len(prices))
	aroonDown := make([]float64, len(prices))

	// Initialize with zero values before we have enough data
	for i := 0; i < period-1; i++ {
		aroonUp[i] = 0
		aroonDown[i] = 0
	}

	// Calculate Aroon indicator for each window
	for i := period - 1; i < len(prices); i++ {
		// Find the index of highest and lowest prices in the lookback period
		highestIndex := findHighestIndex(prices[i-period+1 : i+1])
		lowestIndex := findLowestIndex(prices[i-period+1 : i+1])

		// Calculate Aroon Up: ((period - days since highest) / period) * 100
		daysSinceHigh := period - 1 - highestIndex
		aroonUp[i] = float64(period-daysSinceHigh-1) / float64(period-1) * 100

		// Calculate Aroon Down: ((period - days since lowest) / period) * 100
		daysSinceLow := period - 1 - lowestIndex
		aroonDown[i] = float64(period-daysSinceLow-1) / float64(period-1) * 100
	}

	return aroonUp, aroonDown
}

func findHighestIndex(slice []float64) int {
	maxIndex := 0
	for i := 1; i < len(slice); i++ {
		if slice[i] > slice[maxIndex] {
			maxIndex = i
		}
	}
	return maxIndex
}

// findLowestIndex finds the index of the lowest value in a slice
func findLowestIndex(slice []float64) int {
	minIndex := 0
	for i := 1; i < len(slice); i++ {
		if slice[i] < slice[minIndex] {
			minIndex = i
		}
	}
	return minIndex
}

func calculateUpperBand(sma float64, stdDev float64, multiplier float64) float64 {
	return sma + (stdDev * multiplier)
}

func calculateLowerBand(sma float64, stdDev float64, multiplier float64) float64 {
	return sma - (stdDev * multiplier)
}

func calculateStdDev(prices []float64, period int, mean float64) float64 {
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += math.Pow(prices[i]-mean, 2)
	}
	return math.Sqrt(sum / float64(period))
}
