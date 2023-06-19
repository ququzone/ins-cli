package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ququzone/ins-cli/pkg/ins"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register name",
	Long:  "register ins name",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Must provide name to register")
			return
		}
		if viper.GetString("private-key") == "" {
			fmt.Println("Must provide private key to register")
			return
		}
		if err := ins.Register(
			viper.GetString("rpc"),
			viper.GetString("controller"),
			viper.GetString("resolver"),
			viper.GetString("private-key"),
			viper.GetString("owner"),
			args[0],
		); err != nil {
			fmt.Printf("Register %s fail: %v\n", args[0], err)
		} else {
			fmt.Printf("Register %s successful\n", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().String("rpc", "https://babel-api.mainnet.iotex.io", "blockchain RPC endpoint")
	registerCmd.Flags().String("private-key", "", "private key for regsiter")
	registerCmd.Flags().String("controller", "0x8aA6acF9BFeEE0243578305706766065180E68d4", "IOTXRegistrarController address")
	registerCmd.Flags().String("resolver", "0x41B9132D4661E016A09a61B314a1DFc0038CE3e8", "Public resolver address")
	registerCmd.Flags().String("owner", "", "name owner, default to registrant")
	viper.BindPFlag("rpc", registerCmd.Flags().Lookup("rpc"))
	viper.BindPFlag("private-key", registerCmd.Flags().Lookup("private-key"))
	viper.BindPFlag("controller", registerCmd.Flags().Lookup("controller"))
	viper.BindPFlag("resolver", registerCmd.Flags().Lookup("resolver"))
	viper.BindPFlag("owner", registerCmd.Flags().Lookup("owner"))
}
