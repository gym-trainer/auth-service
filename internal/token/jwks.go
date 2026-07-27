package token

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

func GenerateJWKS(pubKey *rsa.PublicKey) map[string]interface{} {
	eBytes := big.NewInt(int64(pubKey.E)).Bytes()

	return map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"kid": "1",
				"n":   base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes()),
				"e":   base64.RawURLEncoding.EncodeToString(eBytes),
			},
		},
	}
}
