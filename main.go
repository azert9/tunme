package main

import (
	"fmt"
	"github.com/azert9/tunme/pkg/cmd/tunme_cat"
	"github.com/azert9/tunme/pkg/cmd/tunme_relay"
	"github.com/azert9/tunme/pkg/cmd/tunme_tcp"
	"github.com/azert9/tunme/pkg/cmd/tunme_tun"
	"github.com/spf13/cobra"
	"os"
)

func main() {

	// TODO: description, ...
	cmd := cobra.Command{
		Use: os.Args[0],
	}

	tunme_cat.RegisterCmd(&cmd)
	tunme_relay.RegisterCmd(&cmd)
	tunme_tcp.RegisterCmd(&cmd)
	tunme_tun.RegisterCmd(&cmd)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
