package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Downloader = &cobra.Command{
	Use:   "downloader [global options] command [command options] [arguments...]",
	Short: "downloader is a concurrent file downloader",
	Long:  ``,
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

var url string

func init() {
	Downloader.AddCommand()
	Downloader.Flags().StringVarP(&url, "url", "", "", "File URL(required)")
}
