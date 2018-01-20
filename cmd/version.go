package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GITCOMMIT = "HEAD"

	versionCommand = &cobra.Command{
		Use:   "version",
		Short: "Prints the version of hyperledger-fabric-orderer",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("koki/hyperledger-fabric-orderer: %s\n", GITCOMMIT)
		},
	}
)
