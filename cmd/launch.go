package cmd

import (
	"github.com/scbizu/Astral/telegram"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//LaunchCmd impl Launch Command
var LaunchCmd = &cobra.Command{
	Use:   "astral",
	Short: "Launch astral",
	Long:  "Launch command",
	Run: func(cmd *cobra.Command, args []string) {

		if err := telegram.ListenWebHook(true); err != nil {
			logrus.Fatal(err)
		}

		return
	},
}

//Execute exec astral
func Execute() (err error) {
	err = LaunchCmd.Execute()
	return
}
