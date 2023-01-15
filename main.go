package main

import (
	"fmt"
	"github.com/azert9/tunme/pkg/cmd/tunme_cat"
	"github.com/azert9/tunme/pkg/cmd/tunme_relay"
	"github.com/azert9/tunme/pkg/cmd/tunme_tun"
	"github.com/spf13/cobra"
	"os"
)

func main() {

	// TODO: description, ...
	cmd := cobra.Command{
		Use: os.Args[0],
	}

	cmd.AddCommand(&tunme_cat.CobraCmd)
	cmd.AddCommand(&tunme_relay.CobraCmd)
	cmd.AddCommand(&tunme_tun.CobraCmd)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
