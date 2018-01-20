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

	// address at which the GRPC server should listen
	listenAddr string

	// port at which the GRPC server should listen
	listenPort int

	// path to server cert
	serverCertPath string

	// path to server key
	serverKeyPath string

	// path to server root cert cas
	serverRootCAsPath []string

	// path to clietn root cert cas
	clientRootCAsPath []string

	// path to sign-cert
	signCertPath string

	// path to key-store
	keyStorePath string

	// hash function
	hashFunction string

	// MSP Identifier
	mspId string
)

func init() {
	OrdererCommand.Flags().StringVarP(&listenAddr, "addr", "a", "0.0.0.0", "The address at which the orderer service should bind and listen")
	OrdererCommand.Flags().IntVarP(&listenPort, "port", "p", 7050, "The port at which the orderer service should bind and listen")

	OrdererCommand.Flags().StringVarP(&serverCertPath, "server-cert", "", "", "The path to the Server Cert for TLS verification")
	OrdererCommand.Flags().StringVarP(&serverKeyPath, "server-key", "", "", "The path to the Server Key for TLS signing")
	OrdererCommand.Flags().StringSliceVarP(&serverRootCAsPath, "server-root-ca", "", nil, "The path to server root certificate authorities")
	OrdererCommand.Flags().StringSliceVarP(&clientRootCAsPath, "client-root-ca", "", nil, "The path to the client root certificate authorities")

	OrdererCommand.Flags().StringVarP(&signCertPath, "sign-cert", "", "", "The path to sign-cert (public part of the MSP identity)")
	OrdererCommand.Flags().StringVarP(&keyStorePath, "key-store", "", "", "The path to key-store (private key)")
	OrdererCommand.Flags().StringVarP(&hashFunction, "hash-function", "", "SHA256", "The hash function used for identity identifier computation")
	OrdererCommand.Flags().StringVarP(&mspId, "msp-id", "i", "", "ID of the MSP")

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
