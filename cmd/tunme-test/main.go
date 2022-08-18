package main

import (
	"os"
	"tunme/pkg/cmd/tunme_functest"
)

func main() {
	tunme_functest.Main(os.Args[0], os.Args[1:])
}
