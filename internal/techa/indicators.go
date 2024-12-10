package techa

type Indicators struct {
	Trends interface {
		SMA(prices []float64, period int) ([]float64, error)
		EMA(prices []float64, period int) []float64
		DEMA(prices []float64, period int) []float64
		TREMA(prices []float64, period int) []float64
		MACD(prices []float64, fast int, slow int, signal int) ([]float64, []float64, []float64)
		SuperTrend(high, low, close []float64, period int, multiplier float64) []SuperTrendResult
		AvgTrueRange(trValues []float64, period int) []float64
		TrueRange(high, low, close []float64) []float64
		TRIX(prices []float64, period int) []float64
		Aroon(period int, prices []float64) ([]float64, []float64)
	}
	Volatility interface {
		RSI(prices []float64, period int) []float64
		StochRSI(prices []float64, rsiPeriod, stochPeriod int) (float64, error)
		BollingerBands(prices []float64, period int, multiplier float64) ([]float64, []float64, []float64)
	}
	Momentum interface {
		WilliamsR(prices []float64, period int) []float64
	}
}

func NewIndicators() *Indicators {
	return &Indicators{
		Trends:     &Trends{},
		Volatility: &Volatility{},
		Momentum:   &Momentum{},
	}
}
