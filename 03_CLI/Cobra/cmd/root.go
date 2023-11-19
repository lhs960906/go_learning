package cmd

import (
	"fmt"
	"os"

	"github.com/lhs960906/Go-Learning/03_CLI/Cobra/cmd/version"
	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "rootCmd",
	Short: "rootCmd is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				love by spf13 and friends in Go.
				Complete documentation is available at http://hugo.spf13.com`,
	// 执行 rootCmd 的相关 action
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rootCmd执行成功")
	},
	// 当某函数返回 error 不为空时,
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	if err := someFunc(); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(version.GetVersionCmd())
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}
