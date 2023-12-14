package cmd

import (
	"fmt"

	"github.com/lhs960906/Go-Learning/04_Web/Web/01_http/cmd/version"
	"github.com/spf13/cobra"
)

var Downloader = &cobra.Command{
	Use:   "downloader [global options] command [command options] [arguments...]",
	Short: "downloader is a concurrent file downloader",
	Long:  `downloader is a concurrent file downloader`,
	// 执行 downloader 的相关 action
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

var (
	url string
)

func init() {
	Downloader.AddCommand(version.GetVersionCmd())
	Downloader.Flags().StringVarP(&url, "url", "", "", "File URL(required)")
	Downloader.MarkFlagRequired("url")
}
