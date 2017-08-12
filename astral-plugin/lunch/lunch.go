package lunch

import (
	"log"
	"strconv"
	"strings"

	"github.com/scbizu/wechat-go/wxweb"
)

const (
	receiver = ""
	//SPACE defines 4 space globally
	SPACE = "    "
	//ENTER defines enter globally
	ENTER = "\n"
	//DefaultResp is wtf?
	DefaultResp = "WTF!"
)

var (
	lunching, taking bool

	// membersCount int32
	//fromUserID => userinfo
	members = make(map[string]*wxweb.User)

	owner string
	//DataFilePath is the recipe yaml
	DataFilePath = "./astral-plugin/lunch/recipe.yaml"
)

//Register regist the lunch plugin to the bot
func Register(session *wxweb.Session, lunchFunc func(session *wxweb.Session, msg *wxweb.ReceivedMessage)) {
	if lunchFunc == nil {
		lunchFunc = defaultLunch
	}
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(lunchFunc), "callon")
	if err := session.HandlerRegister.EnableByName("callon"); err != nil {
		log.Println(err)
	}

}

func defaultLunch(session *wxweb.Session, msg *wxweb.ReceivedMessage) {

	if !msg.IsGroup {
		return
	}

	var contact *wxweb.User
	orderList := []int{}
	from := session.Bot.UserName
	to := msg.FromUserName

	if to == from {
		contact = session.Cm.GetContactByUserName(msg.ToUserName)
	} else {
		contact = session.Cm.GetContactByUserName(to)
	}
	if contact == nil {
		return
	}
	mm, err := wxweb.CreateMemberManagerFromGroupContact(session, contact)
	if err != nil {
		log.Println(err)
		return
	}
	whoCall := mm.GetContactByUserName(msg.Who)
	msg.Content = strings.TrimSpace(msg.Content)
	if strings.HasPrefix(msg.Content, "/") {
		//Astral的饭团模式开启
		rawCommand := strings.TrimPrefix(msg.Content, "/")
		switch rawCommand {

		case "lunch":
			if lunching {
				session.SendText("lunching,if it was already ended,plz type /end", from, to)
			} else {
				lunching = true
				session.SendText(whoCall.NickName+" created a new lunch,type /join to join in.And type /take to be the lunch owner", from, to)
			}
		case "take":
			if !lunching {
				session.SendText("type /new to new lunch first", from, to)
			} else if owner != "" {
				session.SendText("sorry,the owner is "+owner, from, to)
			} else {
				taking = true
				owner = whoCall.UserName
				session.SendText(whoCall.NickName+" now is the luncher owner.Type `/random ?`OR `/manual (?,?,?)`(type /baocan to check the index of the menu) to order, e.g: `/random 10` OR `/manual (1,2,3)`", from, to)
			}

		case "baocan":
			orderListStr := getAllRecipeInfo()
			session.SendText(orderListStr, from, to)
		case "send":
			if taking && to == owner {
				session.SendText("您好,我们点菜,12点左右过来吃。菜单是:\n "+showRecipe(orderList)+"\n 麻烦了,O(∩_∩)O谢谢~", from, receiver)
				session.SendText("recipe was sent,type /star ?(0 - 5) to star this recipe after lunch.", from, to)
				//reset
				taking = false
			} else {
				session.SendText(DefaultResp, from, to)
			}

		case "join":
			if !lunching {
				session.SendText("no lunch now,type /lunch to create a new lunch", from, to)
			} else {
				if _, ok := members[whoCall.UserName]; ok {
					session.SendText(DefaultResp, from, to)
				}
				members[whoCall.UserName] = whoCall
				membersCount := len(members)
				session.SendText(whoCall.NickName+" joined in,now has "+strconv.Itoa(int(membersCount))+" members", from, to)
			}
		case "end":
			if !lunching {
				session.SendText("already ended", from, to)
			} else {
				session.SendText("lunch ended", from, to)
				lunching = false
				resetMember()
			}
		case "quit":
			if _, ok := members[whoCall.UserName]; !ok {
				session.SendText(DefaultResp, from, to)
			} else {
				delete(members, whoCall.UserName)
				session.SendText(whoCall.NickName+" quited,now has "+strconv.Itoa(int(len(members)))+" members", from, to)
			}
		case "whoami":
			session.SendText("FromUserName: "+to+ENTER+"Username: "+whoCall.UserName+ENTER+"Nick Name: "+whoCall.NickName+ENTER+"Uin: "+strconv.Itoa(whoCall.Uin), from, to)
		case "members":
			if lunching {
				respStr := "members in lunch pool:" + ENTER
				membersSlice := []string{}
				for _, v := range members {
					membersSlice = append(membersSlice, v.NickName)
				}
				membersStr := strings.Join(membersSlice, ",")
				respStr += "members: " + membersStr + ENTER
				if _, ok := members[owner]; ok {
					respStr += "owner: " + members[owner].NickName
				}
				session.SendText(respStr, from, to)
			} else {
				session.SendText("no lunch", from, to)
			}
		default:
			if strings.Contains(rawCommand, "manual") {
				if taking && to == owner {
					// println(rawCommand)
					params := strings.TrimPrefix(rawCommand, "manual(")
					params = strings.TrimSuffix(params, ")")
					println(params)
					orders := strings.Split(params, ",")
					for _, o := range orders {
						oi, err := strconv.ParseInt(o, 10, 64)
						if err != nil {
							log.Println(err)
							session.SendText(DefaultResp, from, to)
						}
						orderList = append(orderList, int(oi))
					}
					orderListStr := showRecipe(orderList)
					session.SendText(orderListStr, from, to)
				} else {
					session.SendText(DefaultResp, from, to)
				}
			} else if strings.Contains(rawCommand, "random") {
				if taking && to == owner {
					params := strings.TrimSpace(strings.TrimPrefix(rawCommand, "random"))
					count, err := strconv.ParseInt(params, 10, 64)
					if err != nil {
						log.Println(err)
						session.SendText(DefaultResp, from, to)
					}
					orderListsName := getRecipe(int(count))
					if len(orderListsName) > 0 {
						session.SendText(showRecipeByName(orderListsName)+"\n type /send to send this order to BAOCAN", from, to)
					}
					session.SendText("plz reduce request recipe counts", from, to)
				} else {
					session.SendText(DefaultResp, from, to)
				}
			} else if strings.Contains(rawCommand, "star") {
				if lunching && to == owner {
					star := strings.TrimSpace(strings.TrimPrefix(rawCommand, "star"))
					starInt, err := strconv.ParseInt(star, 10, 64)
					if err != nil {
						log.Println(err)
						session.SendText(DefaultResp, from, to)
					}
					err = starRecipe(orderList, int(starInt))
					if err != nil {
						log.Println(err)
						session.SendText(DefaultResp, from, to)
					}
					session.SendText("The recipe is "+star+" ✨,it will work next time.type /end to end.", from, to)
					owner = ""
				} else {
					session.SendText(DefaultResp, from, to)
				}
			} else if strings.Contains(rawCommand, "pin") {
				var PIN string
				pinMsg := strings.TrimSpace(strings.TrimPrefix(rawCommand, "pin"))
				for _, m := range mm.Group.MemberList {
					if m.UserName != whoCall.UserName {
						PIN += "@" + m.NickName + " "
					}
				}
				PIN += pinMsg
				session.SendText(PIN, from, to)
			} else {
				session.SendText(DefaultResp, from, to)
			}
		}
	}
}

func resetMember() {
	for k := range members {
		delete(members, k)
	}
}
