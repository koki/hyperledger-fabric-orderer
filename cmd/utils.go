package cmd

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/hyperledger/fabric/core/comm"
	"github.com/spf13/cobra"

	"github.com/koki/hyperledger-fabric-orderer/config"
)

func validateOrdererFlags(c *cobra.Command, args []string) error {
	// validate the ip address
	if ip := net.ParseIP(listenAddr); ip == nil {
		return fmt.Errorf("Invalid address %s", listenAddr)
	}

	// validate the port
	if listenPort < 0 || listenPort > 65535 {
		return fmt.Errorf("Invalid port %d", listenPort)
	}

	return nil
}

func extractOrdererConfig(c *cobra.Command, args []string) (*config.OrdererConfig, error) {
	if c == nil {
		return nil, nil
	}

	var err error

	useTLS := false
	requireClientCert := false

	var serverCert []byte
	var serverKey []byte
	var serverRootCAs [][]byte
	var clientRootCAs [][]byte

	var signCert []byte
	var keyStore []byte

	if serverCertPath != "" {
		serverCert, err = ioutil.ReadFile(serverCertPath)
		if err != nil {
			return nil, err
		}
		useTLS = true
	}

	if serverKeyPath != "" {
		serverKey, err = ioutil.ReadFile(serverKeyPath)
		if err != nil {
			return nil, err
		}
	}

	if serverRootCAsPath != nil {
		for i := range serverRootCAsPath {
			serverRootCAPath := serverRootCAsPath[i]
			serverRootCA, err := ioutil.ReadFile(serverRootCAPath)
			if err != nil {
				return nil, err
			}
			serverRootCAs = append(serverRootCAs, serverRootCA)
		}
	}

	if clientRootCAsPath != nil {
		// mis-configuration of non-nil and empty
		requireClientCert = true
		for i := range clientRootCAsPath {
			clientRootCAPath := clientRootCAsPath[i]
			clientRootCA, err := ioutil.ReadFile(clientRootCAPath)
			if err != nil {
				return nil, err
			}
			clientRootCAs = append(clientRootCAs, clientRootCA)
		}
	}

	if signCertPath != "" {
		signCert, err = ioutil.ReadFile(signCertPath)
		if err != nil {
			return nil, err
		}
	}

	if keyStorePath != "" {
		keyStore, err = ioutil.ReadFile(keyStorePath)
		if err != nil {
			return nil, err
		}
	}

	var typedHashFunction config.IdentityIdentifierHashFunction
	switch hashFunction {
	case "SHA256":
		typedHashFunction = config.SHA256
	case "SHA384":
		typedHashFunction = config.SHA384
	case "SHA3_256":
		typedHashFunction = config.SHA3_256
	case "SHA3_384":
		typedHashFunction = config.SHA3_384
	default:
		return nil, fmt.Errorf("unsupported hash function %s", hashFunction)
	}

	secureConfig := comm.SecureServerConfig{
		UseTLS:            useTLS,
		RequireClientCert: requireClientCert,
		ServerCertificate: serverCert,
		ServerKey:         serverKey,
		ServerRootCAs:     serverRootCAs,
		ClientRootCAs:     clientRootCAs,
	}

	// Note: do not change order
	config := &config.OrdererConfig{
		secureConfig,
		listenAddr,
		listenPort,
		keyStore,
		signCert,
		typedHashFunction,
		mspId,
	}

	return config, err
}
