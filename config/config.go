package config

import (
	"github.com/hyperledger/fabric/core/comm"
)

// Note: If you change order here, update cmd/utils.go
type OrdererConfig struct {
	comm.SecureServerConfig

	Address string
	Port    int

	KeyStore []byte
	SignCert []byte

	IdentityIdentifierHashFunction IdentityIdentifierHashFunction

	MspId string
}

type IdentityIdentifierHashFunction string

const (
	SHA256   IdentityIdentifierHashFunction = "SHA256"
	SHA384   IdentityIdentifierHashFunction = "SHA384"
	SHA3_256 IdentityIdentifierHashFunction = "SHA_256"
	SHA3_384 IdentityIdentifierHashFunction = "SHA_384"
)
