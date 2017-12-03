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
		session, err := wxweb.CreateSession(nil, nil, wxweb.WEB_MODE)
		if err != nil {
			log.Fatal(err)
			return
		}

		if isHTTP {
			go http.ListenAndServe(":8080", http.FileServer(http.Dir("./")))
		}

		plugin.RegisterAllEnabledPlugins(session)

		go func() {
			if err := telegram.PullAndReply(); err != nil {
				log.Fatal(err)
			}
		}()

		if err := session.LoginAndServe(false); err != nil {
			log.Fatal(err)
		}
		// for {
		// 	if err := session.LoginAndServe(false); err != nil {
		// 		logs.Error("session exit, %s", err)
		// 		for i := 0; i < 3; i++ {
		// 			logs.Info("trying re-login with cache")
		// 			if err := session.LoginAndServe(true); err != nil {
		// 				logs.Error("re-login error, %s", err)
		// 			}
		// 			time.Sleep(3 * time.Second)
		// 		}
		// 		if session, err = wxweb.CreateSession(nil, session.HandlerRegister, wxweb.TERMINAL_MODE); err != nil {
		// 			logs.Error("create new sesion failed, %s", err)
		// 			break
		// 		}
		// 	} else {
		// 		logs.Info("closed by user")
		// 		break
		// 	}
		// }
		return
	},
}

//Execute exec astral
func Execute() (err error) {
	err = LaunchCmd.Execute()
	return
}
