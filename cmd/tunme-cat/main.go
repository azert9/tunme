package main

import (
	"os"
	"tunme/pkg/cmd/tunme_cat"
)

func main() {
	tunme_cat.Main(os.Args[0], os.Args[1:])
}
