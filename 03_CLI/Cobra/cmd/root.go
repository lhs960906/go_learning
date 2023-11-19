package Cobra

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				love by spf13 and friends in Go.
				Complete documentation is available at http://hugo.spf13.com`,
	// 如果有相关的 action 要执行，请取消下面这行代码的注释
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("测试命令执行")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("test")
}
