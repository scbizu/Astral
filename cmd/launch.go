package cmd

import (
	"github.com/scbizu/Astral/internal/discord"
	"github.com/scbizu/Astral/internal/telegram"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AstralCmd defines astral main cmd
var AstralCmd = &cobra.Command{
	Use:   "astral",
	Short: "Astral Bot",
	Run: func(_ *cobra.Command, _ []string) {
		bot, err := telegram.NewBot(true)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		if err := bot.ServeBotUpdateMessage(); err != nil {
			logrus.Fatal(err)
			return
		}
	},
}

// ServiceCmd defines astral service cmd
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Astral is a self-host bot of @scnace",
}

// TGService defines astral bot command service
var TGService = &cobra.Command{
	Use:   "telegram",
	Short: "a telegram bot",
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
	},
}

var DiscordService = &cobra.Command{
	Use:   "discord",
	Short: "a discord bot",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := discord.NewBot()
		if err != nil {
			logrus.Fatal(err)
			return err
		}
		return nil
	},
}

// Execute exec astral
func Execute() (err error) {
	ServiceCmd.AddCommand(TGService)
	ServiceCmd.AddCommand(DiscordService)
	AstralCmd.AddCommand(ServiceCmd)
	return AstralCmd.Execute()
}
