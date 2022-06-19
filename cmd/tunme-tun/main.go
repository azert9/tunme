package main

import (
	"os"
	"tunme/pkg/cmd/tunme_tun"
)

func main() {
	tunme_tun.Main(os.Args[0], os.Args[1:])
}
