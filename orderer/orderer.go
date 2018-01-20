package orderer

import (
	"fmt"
	"net"

	"github.com/hyperledger/fabric/core/comm"

	"github.com/koki/hyperledger-fabric-orderer/config"
)

const (
	TCP = "tcp"
)

func Run(config *config.OrdererConfig) error {
	if config == nil {
		return fmt.Errorf("nil config obtained")
	}

	signer, err := NewSignerFromConfig(config)
	if err != nil {
		return err
	}

	signer.Sign(nil)

	addr := fmt.Sprintf("%s:%d", config.Address, config.Port)
	listener, err := net.Listen(TCP, addr)
	if err != nil {
		return err
	}

	server, err := comm.NewGRPCServerFromListener(listener, config.SecureServerConfig)
	if err != nil {
		return err
	}

	return server.Start()
}
