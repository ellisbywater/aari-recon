package coinbase

import (
	coinbase "aari-recon/internal/coinbase/credentials"
	"aari-recon/internal/env"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const (
	jwtRequestHost   = "api.coinbase.com"
	jwtRequestPath   = "/api/v3/brokerage/accounts"
	jwtRequestMethod = "GET"
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
	keySecret, err := env.GetStringNoFallback("COINBASE_PRIVATE_KEY")
	fmt.Println("COINBASE PRIVATE KEY >>>> ", keySecret)
	if err != nil {
		return "", fmt.Errorf("jwt: Could not decode")
	}
	block, _ := pem.Decode([]byte(coinbase.COINBASE_PRIVATE_KEY))
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
