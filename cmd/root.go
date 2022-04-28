package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	version bool
)

var rootCmd = &cobra.Command{
	Use:   "restful-api",
	Short: "restful-api crud demo",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if version {
			fmt.Println("0.0.1")
			return nil
		}
		return errors.New("no flags find")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "the cloud-station-cli version")
}
