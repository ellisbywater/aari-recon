package techa

import "math"

type SMA struct {
	values []float64
	window int
}

func NewSMA(data []float64, window int) *SMA {
	sma := &SMA{values: make([]float64, 0), window: window}
	for i := 0; i < window; i++ {
		sum := 0.0
		for j := i; j < i+window; j++ {
			sum += data[j]
		}
		sma.values = append(sma.values, sum/float64(window))
	}
	return sma
}

func (sma *SMA) Add(value float64) {
	sum := 0.0
	for i := 1; i < len(sma.values); i++ {
		sum += sma.values[i-1]
	}
	sma.values = append([]float64{value}, sma.values...)

	// calculate the new average
	sum += value
	sma.values[0] = sum / float64(len(sma.values))
}

type EMA struct {
	alpha float64
	sum   float64
	count int
}

func NewEMA(alpha float64) *EMA {
	return &EMA{alpha: alpha, sum: 0, count: 0}
}

func (ema *EMA) Add(value float64) {
	ema.sum += (value - ema.sum/float64(ema.count)) * ema.alpha
	ema.count++
}

func (ema *EMA) Value() float64 {
	return ema.sum / float64(ema.count)
}

type MACD struct {
	fast         int
	slow         int
	signal       int
	values       []float64
	signalValues []float64
	delta        []float64
}

func NewMACD(fast int, slow int, signal int) *MACD {
	return &MACD{fast: fast, slow: slow, signal: signal}
}

func (macd *MACD) Calculate(prices []float64) ([]float64, []float64, []float64) {
	shortEMA := make([]float64, len(prices))
	longEMA := make([]float64, len(prices))
	macdvalue := make([]float64, len(prices))
	signal := make([]float64, len(prices))
	delta := make([]float64, len(prices))

	shortEMA[0] = prices[0]
	longEMA[0] = prices[0]

	for i := 1; i < len(prices); i++ {
		shortEMA[i] = (prices[i]*2 + shortEMA[i-1]*(float64(macd.fast)-2)) / float64(macd.fast)
		longEMA[i] = (prices[i]*2 + longEMA[i-1]*(float64(macd.slow)-2)) / float64(macd.slow)

		// Calculate MACD line
		macdvalue[i] = longEMA[i] - shortEMA[i]

		// Calculate Signal line
		signal[i] = macdvalue[i]*2 + signal[i-1]*(float64(macd.signal)-2)/float64(macd.signal)

		// Calculate Convergence/Divergence
		delta[i] = macdvalue[i] - signal[i]
	}

	macd.values = macdvalue
	macd.signalValues = signal
	macd.delta = delta

	return macdvalue, signal, delta
}

type VariableEWMA struct {
	alpha  float64
	values []float64
}

func NewVariableEWMA(alpha float64) *VariableEWMA {
	return &VariableEWMA{alpha: alpha}
}

func (e *VariableEWMA) Calculate(prices []float64) []float64 {
	e.values = make([]float64, len(prices))
	for i, price := range prices {
		if i < 10 {
			e.values = append(e.values, price)
		} else {
			newVal := e.values[i-1]*e.alpha + price*(1-e.alpha)
			e.values = append(e.values, newVal)
		}
	}
	return e.values
}

type BollingerBands struct {
	period     int
	multiplier float64
	SMA        []float64
	StdDev     []float64
	Upper      []float64
	Lower      []float64
}

func NewBollingerBands(period int, multiplier float64) *BollingerBands {
	return &BollingerBands{period: period, multiplier: multiplier}
}

func (bb *BollingerBands) Calculate(data []float64) ([]float64, []float64, []float64) {
	smaVals := make([]float64, len(data))
	stdVals := make([]float64, len(data))
	upperVals := make([]float64, len(data))
	lowerVals := make([]float64, len(data))

	for i := 0; i < len(data); i++ {
		if i >= bb.period-1 {
			smaVals[i] = data[i]
			stdVals[i] = data[i]
			upperVals[i] = data[i]
			lowerVals[i] = data[i]
		} else {
			smaVals[i] = CalculateSMA(data[i:i-bb.period], bb.period)
			stdVals[i] = calculateStdDev(data[i:i-bb.period], bb.period, smaVals[i])
			upperVals[i] = calculateUpperBand(smaVals[i], stdVals[i], bb.multiplier)
			lowerVals[i] = calculateLowerBand(smaVals[i], stdVals[i], bb.multiplier)
		}
	}

	return smaVals, upperVals, lowerVals

}

func CalculateSMA(prices []float64, period int) float64 {
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	return sum / float64(period)
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
