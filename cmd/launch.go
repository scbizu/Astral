package cmd

import (
	"github.com/scbizu/Astral/internal/telegram"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AstralCmd defines astral main cmd
var AstralCmd = &cobra.Command{
	Use:   "astral",
	Short: "Astral Telegram Bot & Push Service",
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
		bot, err := telegram.NewBot(true)
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

// ServerInfoService defines autopush astral server channel message service
var ServerInfoService = &cobra.Command{
	Use:   "AstralServerMessage",
	Short: "Astral server message service",
	Run: func(cmd *cobra.Command, args []string) {
		bot, err := telegram.NewBot(true)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		bot.ServePushAstralServerMessage()
		return
	},
}

// CNSC2EventInfoService defines autopush CNSC2Event channle message service
var CNSC2EventInfoService = &cobra.Command{
	Use:   "SC2EventInfo",
	Short: "CNSC2 Event info service",
	Run: func(cmd *cobra.Command, args []string) {
		bot, err := telegram.NewBot(true)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		bot.ServePushSC2Event()
		return
	},
}

//Execute exec astral
func Execute() (err error) {
	ServiceCmd.AddCommand(CommandService)
	ServiceCmd.AddCommand(ServerInfoService)
	ServiceCmd.AddCommand(CNSC2EventInfoService)
	AstralCmd.AddCommand(ServiceCmd)
	return AstralCmd.Execute()
}
