package cmd

import (
	"github.com/scbizu/Astral/internal/telegram"
	"github.com/scbizu/Astral/service/push"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AstralCmd defines astral main cmd
var AstralCmd = &cobra.Command{
	Use:   "astral",
	Short: "Astral Bot & Push Service",
}

// ServiceCmd defines astral service cmd
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Astral service cmd(command,autopush)",
}

// CommandService defines astral bot command service
var CommandService = &cobra.Command{
	Use:   "Command",
	Short: "Astral command service",
	Run: func(cmd *cobra.Command, args []string) {
		bot, err := telegram.NewBot(false)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		if err := bot.ServeBotUpdateMessage(); err != nil {
			logrus.Fatal(err)
			return
		}

		return
	},
}

// CNSC2EventInfoService defines autopush CNSC2Event channle message service
var CNSC2EventInfoService = &cobra.Command{
	Use:   "SC2EventInfo",
	Short: "CNSC2 Event info service",
	Run: func(cmd *cobra.Command, args []string) {
		s := push.NewPushService()
		if err := s.ServePushSC2Event(); err != nil {
			logrus.Fatal(err)
		}
	},
}

//Execute exec astral
func Execute() (err error) {
	ServiceCmd.AddCommand(CommandService)
	ServiceCmd.AddCommand(CNSC2EventInfoService)
	AstralCmd.AddCommand(ServiceCmd)
	return AstralCmd.Execute()
}
