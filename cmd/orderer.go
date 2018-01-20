package cmd

import (
	"flag"

	_ "github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/koki/hyperledger-fabric-orderer/config"
	"github.com/koki/hyperledger-fabric-orderer/orderer"
)

var (
	// root command of this program
	OrdererCommand = &cobra.Command{
		Use:   "hyperledger-fabric-orderer",
		Short: "Hyperledger Fabric Orderer",
		Long:  "Hyperledger Node for total order broadcast",
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateOrdererFlags(c, args); err != nil {
				return err
			}

			config, err := extractOrdererConfig(c, args)
			if err != nil {
				return err
			}

			return runOrderer(config)
		},
		SilenceUsage: true,
		Example: `
`,
	}
)

func init() {
	// parse the go default flagset to get flags for glog and other packages in future
	OrdererCommand.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// defaulting this to true so that logs are printed to console
	flag.Set("logtostderr", "true")

	// suppress the incorrect prefix in glog output
	flag.CommandLine.Parse([]string{})

	// add sub-commands
	OrdererCommand.AddCommand(versionCommand)
}

func runOrderer(config *config.OrdererConfig) error {
	return orderer.Run(config)
}
