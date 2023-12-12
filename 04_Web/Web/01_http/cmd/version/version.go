package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "downloader version",
	Short: "Print the version number of downloader",
	Long:  `All software has versions. This is Downloader's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Network file downloader v0.9")
	},
}

func GetVersionCmd() *cobra.Command {
	return versionCmd
}
