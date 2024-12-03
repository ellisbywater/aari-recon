package main

import (
	"aari-recon/internal/coinbase"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	OneMinGran     = "60"
	FiveMinGran    = "300"
	FifteenMinGran = "900"
	OneHourGran    = "3600"
	SixHourGran    = "21600"
	OneDayGran     = "86400"
)

type Assumption struct {
	Text      string `json:"text"`
	Sentiment bool   `json:"sentiment"`
}

type Asset struct {
	Symbol      string       `json:"symbol"`
	Assumptions []Assumption `json:"assumptions"`
	Market      string       `json:"market"`
}

type AssetConfig struct {
	CandleInterval   int64 `json:"candle_interval"`
	ResearchInterval int64 `json:"research_interval"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env vars")
		return
	}
	environs := os.Environ()
	for _, s := range environs {
		fmt.Println(s)
	}

	jwt, err := coinbase.BuildJwt()
	if err != nil {
		fmt.Println("Error building jwt:  ", err)
	} else {
		fmt.Printf("JWT BUILD SUCCESS: %s", jwt)
	}

}
