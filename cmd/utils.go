package cmd

import (
	"github.com/koki/hyperledger-fabric-orderer/config"
	"github.com/spf13/cobra"
)

func validateOrdererFlags(c *cobra.Command, args []string) error {
	return nil
}

func extractOrdererConfig(c *cobra.Command, args []string) (*config.OrdererConfig, error) {
	return nil, nil
}
