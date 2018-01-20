package main

import (
	"os"

	"github.com/koki/hyperledger-fabric-orderer/cmd"
)

/*
  koki/hyperledger-fabric-orderer
	-------------------------------

	Hyperledger orderer with better code and design

	Explicit Goals
	xxxxxxxxxxxxxx

	1. Better Code
	2. Better Documentation
	3. Microservices Architecture
	4. Highly Scalable and designed to be orchestrated

	Docs available at: https://docs.koki.io/orderer

*/

func main() {
	if err := cmd.OrdererCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
