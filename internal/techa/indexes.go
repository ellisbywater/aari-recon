package techa

type RSI struct {
	AvgGains  []float64
	AvgLosses []float64
	RS        []float64
	RSIValues []float64
	period    int
}

func NewRSI(period int) *RSI {
	return &RSI{period: period}
}

func (rsi *RSI) Calculate(closes []float64) {
	rsi.AvgGains = make([]float64, len(closes))
	rsi.AvgLosses = make([]float64, len(closes))
	rsi.RS = make([]float64, len(closes))
	rsi.RSIValues = make([]float64, len(closes))

	for i := 0; i < len(closes); i++ {
		if i < rsi.period {
			rsi.AvgGains[i] = 0.0
			rsi.AvgLosses[i] = 0.0
			rsi.RS[i] = 0.0
			rsi.RSIValues[i] = 0.0
		} else {
			avgGain, avgLoss, rs, rsiValue := calculateCurrentRSI(closes[i-rsi.period : i])
			rsi.AvgGains[i] = avgGain
			rsi.AvgLosses[i] = avgLoss
			rsi.RS[i] = rs
			rsi.RSIValues[i] = rsiValue
		}
	}
}

func calculateCurrentRSI(closes []float64) (float64, float64, float64, float64) {
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
	return avgGain, avgLoss, rs, rsiValue
}

type StochRSI struct {
	rsiValues      []float64
	stochRSIValues []float64
	hh             float64
	ll             float64
	period         int
}

func (s *StochRSI) Calculate(closes []float64) {
	highest, lowest := highestAndLowest(closes)
	s.hh = highest
	s.ll = lowest
	delta := s.hh - s.ll
	s.stochRSIValues = make([]float64, len(closes))
	s.rsiValues = make([]float64, len(closes))

	for i, _ := range closes {
		if i < s.period {
			continue
		}
		_, _, _, rsi := calculateCurrentRSI(closes[i-s.period : i])
		s.rsiValues[i] = rsi
		s.stochRSIValues[i] = (rsi - s.ll) / delta * 100
	}
}

func (s *StochRSI) Values() ([]float64, []float64) {
	return s.stochRSIValues, s.rsiValues
}
