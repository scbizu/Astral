package cmd

import (
	"log"
	"net/http"

	"github.com/scbizu/Astral/astral-plugin"
	"github.com/scbizu/Astral/telegram"
	"github.com/scbizu/wechat-go/wxweb"
	"github.com/spf13/cobra"
)

var isHTTP bool

// var towardsTG bool

func init() {
	LaunchCmd.PersistentFlags().BoolVarP(&isHTTP, "http", "p", false, "put the qrcode into the website")
	// LaunchCmd.PersistentFlags().BoolVarP(&towardsTG, "tg", "t", false, "send to telegram")
}

//LaunchCmd impl Launch Command
var LaunchCmd = &cobra.Command{
	Use:   "astral",
	Short: "Launch astral",
	Long:  "Launch command",
	Run: func(cmd *cobra.Command, args []string) {

		// if err := telegram.PullAndReply(); err != nil {
		// 	log.Fatal(err)
		// }

		if err := telegram.ListenWebHook(true); err != nil {
			log.Fatal(err)
		}

		//wechat Launch in a go routine
		go func() {
			session, err := wxweb.CreateSession(nil, nil, wxweb.WEB_MODE)
			if err != nil {
				log.Printf("create wechat session failed:%s", err.Error())
				return
			}

			if isHTTP {
				go http.ListenAndServe(":8080", http.FileServer(http.Dir("./")))
			}

			plugin.RegisterWechatEnabledPlugins(session)

			if err := session.LoginAndServe(false); err != nil {
				log.Printf("wechat listener has an error:%s", err.Error())
			}
		}()

		return
	},
}

//Execute exec astral
func Execute() (err error) {
	err = LaunchCmd.Execute()
	return
}
