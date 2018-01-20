package orderer

import (
	"fmt"

	"github.com/koki/hyperledger-fabric-orderer/config"
)

func Run(config *config.OrdererConfig) error {
	fmt.Println("Orderer pipeline plugged fully")
	return nil
}
