package main

import (
	"log"

	"github.com/scbizu/wechat-go/wxweb"
	"github.com/scbizu/wxgo/astral-plugin/lunch"
)

func main() {
	session, err := wxweb.CreateSession(nil, nil, wxweb.TERMINAL_MODE)
	if err != nil {
		log.Fatal(err)
		return
	}
	// replier.Register(session, autoReply)
	lunch.Register(session, nil)

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
}

// func autoReply(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
// 	if !msg.IsGroup {
// 		session.SendText("留言收到了,🐔正在认真搬砖哦~", session.Bot.UserName, msg.FromUserName)
// 	}
// }
