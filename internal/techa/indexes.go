package techa

type RSI struct {
	AvgGain  float64
	AvgLoss  float64
	RS       float64
	RSIValue float64
	period   int
}

func NewRSI(period int) *RSI {
	return &RSI{period: period}
}

func (rsi *RSI) Calculate(closes []float64) {
	var (
		gainSum   float64
		lossSum   float64
		gainCount int
		lossCount int
		avgGain   float64
		avgLoss   float64
		rs        float64
		rsiValue  float64
	)

	for i := 0; i < len(closes); i++ {
		if i < rsi.period {
			continue
		}
		var gain, loss float64
		if closes[i] > closes[i-1] {
			gain = closes[i] - closes[i-1]
			gainCount++
			gainSum += gain
		} else if closes[i] < closes[i-1] {
			loss = closes[i-1] - closes[i]
			lossCount++
			lossSum += loss
		}

		if gainCount > 0 && lossCount > 0 {
			avgGain = gainSum / float64(gainCount)
			avgLoss = lossSum / float64(lossCount)
			rs = avgGain / avgLoss
			rsiValue = 100 - (100 / (1 + rs))
		}
	}

	rsi.AvgGain = avgGain
	rsi.AvgLoss = avgLoss
	rsi.RS = rs
	rsi.RSIValue = rsiValue
}

type StochRSI struct {
	rsiValues []float64
	hh        float64
	ll        float64
}

func (s *StochRSI) Calculate(prices []float64) []float64 {
	highest, lowest := highestAndLowest(prices)
	s.hh = highest
	s.ll = lowest
	delta := s.hh - s.ll
	stochrsiValues := make([]float64, len(prices))

	for i, rsi := range s.rsiValues {
		stochrsiValues[i] = (rsi - s.ll) / delta * 100
	}
	return stochrsiValues
}
