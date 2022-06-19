package main

import (
	"os"
	"tunme/pkg/cmd/tunme_relay"
)

func main() {
	tunme_relay.Main(os.Args[0], os.Args[1:])
}
