package coinbase

import (
	"aari-recon/internal/env"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const (
	jwtRequestHost   = "api.coinbase.com"
	jwtRequestPath   = "/api/v3/brokerage/accounts"
	jwtRequestMethod = "GET"
	OneMinGran       = "ONE_MINUTE"
	FiveMinGran      = "FIVE_MINUTE"
	FifteenMinGran   = "FIFTEEN_MINUTE"
	ThirtyMinGran    = "THIRTY_MINUTE"
	OneHour          = "ONE_HOUR"
	TwoHour          = "TWO_HOUR"
	SixHourGran      = "SIX_HOUR"
	OneDayGran       = "ONE_DAY"
)

var max = big.NewInt(math.MaxInt64)

type nonceSource struct{}

func (n nonceSource) Nonce() (string, error) {
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return r.String(), nil
}

type APIKeyClaims struct {
	*jwt.Claims
	URI string `json:"uri"`
}

func BuildJwt() (string, error) {
	uri := fmt.Sprintf("%s %s%s", jwtRequestMethod, jwtRequestHost, jwtRequestPath)

	privateKey, err := env.GetStringNoFallback("COINBASE_PRIVATE_KEY")
	if err != nil {
		return "", fmt.Errorf("jwt: error fetching private key")
	}
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", fmt.Errorf("jwt: Could not decode private key")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}
	keyName, err := env.GetStringNoFallback("COINBASE_KEY_NAME")
	if err != nil {
		return "", fmt.Errorf("no key name")
	}
	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.ES256, Key: key},
		(&jose.SignerOptions{NonceSource: nonceSource{}}).WithType("JWT").WithHeader("kid", keyName),
	)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}
	cl := &APIKeyClaims{
		Claims: &jwt.Claims{
			Subject:   keyName,
			Issuer:    "cdp",
			NotBefore: jwt.NewNumericDate(time.Now()),
			Expiry:    jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
		},
		URI: uri,
	}
	jwtString, err := jwt.Signed(sig).Claims(cl).Serialize()
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}
	return jwtString, nil
}

func FetchAssetCandles(ticker string, start string, end string, granularity string) error {
	jwt, err := BuildJwt()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/products/%s/candles?start=%s&end=%s&granularity=%s",
			ticker,
			start,
			end,
			granularity),
		nil,
	)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", jwt)},
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("response >> ", res)
	return nil
}

func FetchAsset(ticker string) error {
	jwt, err := BuildJwt()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/products/%s", ticker),
		nil,
	)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", jwt)},
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("response >> ", res)
	return nil
}
