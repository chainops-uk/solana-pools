package main

import (
	"fmt"
	"github.com/dfuse-io/solana-go/rpc"
	"github.com/everstake/solana-pools/pkg/pools/parrot"
	"github.com/spf13/cobra"
	"os"
)

func main()  {
	var command = &cobra.Command{
		Use:   "solana pools",
		Short: "start solana pools application",
		Long:  `start solana pools application`,
		RunE: func(cmd *cobra.Command, args []string) error {
			rpcCli := rpc.NewClient("https://api.mainnet-beta.solana.com")
			pool := parrot.New(rpcCli)
			poolData, err  := pool.GetData("AMjGNE12gNoZnrU68AGxUibYEjrGPgpPk3EYG5MZCiZQ")
			if err != nil {
				return fmt.Errorf("pool.GetData: %s", err.Error())
			}
			fmt.Printf("%+v \n", poolData)
			return nil
		},
	}

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}