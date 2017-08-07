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

)

var (
	lunching, taking bool

	membersCount int32
	//fromUserID => userinfo
	members map[string]*wxweb.User

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
	if msg.FromUserName == session.Bot.UserName {
		contact = session.Cm.GetContactByUserName(msg.ToUserName)
	} else {
		contact = session.Cm.GetContactByUserName(msg.FromUserName)
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
				session.SendText("lauching,if it was already ended,plz type /end", session.Bot.UserName, msg.FromUserName)
			} else {
				lunching = true
				session.SendText(whoCall.NickName+" created a new lunch,type /join to join in.And type /take to be the lunch owner", session.Bot.UserName, msg.FromUserName)
			}
		case "take":
			if !lunching {
				session.SendText("type /new to new lunch first", session.Bot.UserName, msg.FromUserName)
			} else if owner != "" {
				session.SendText("sorry,the owner is "+owner, session.Bot.UserName, msg.FromUserName)
			} else {
				taking = true
				owner = msg.FromUserName
				session.SendText(whoCall.NickName+" now is the luncher owner.Type `/random ?`OR `/manual (?,?,?)`(type /baocan to check the index of the menu) to order, e.g: `/random 10` OR `/manual (1,2,3)`", session.Bot.UserName, msg.FromUserName)
			}

		case "baocan":
			orderListStr := getAllRecipeInfo()
			session.SendText(orderListStr, session.Bot.UserName, msg.FromUserName)
		case "send":
			if taking && msg.FromUserName == owner {
				session.SendText("您好,我们点菜,12点左右过来吃。菜单是:\n "+showRecipe(orderList)+"\n 麻烦了,O(∩_∩)O谢谢~", session.Bot.UserName, receiver)
				session.SendText("recipe was sent,type /star ?(0 - 5) to star this recipe after lunch.", session.Bot.UserName, msg.FromUserName)
				//reset
				taking = false
			} else {
				session.SendText("WTF!Dont kidding me!", session.Bot.UserName, msg.FromUserName)
			}
		case "star":

		case "join":
			if !lunching {
				session.SendText("no lunch now,type /lunch to create a new lunch", session.Bot.UserName, msg.FromUserName)
			} else {
				// membersCount = atomic.AddInt32(&membersCount, 1)
				membersCount = membersCount + 1
				session.SendText(whoCall.NickName+" joined in,now has "+strconv.Itoa(int(membersCount))+" members", session.Bot.UserName, msg.FromUserName)
			}
		case "end":
			if !lunching {
				session.SendText("already ended", session.Bot.UserName, msg.FromUserName)
			} else {
				session.SendText("lunch ended", session.Bot.UserName, msg.FromUserName)
				lunching = false
				membersCount = 0
			}
		default:
			if strings.Contains(rawCommand, "manual") {
				if taking && msg.FromUserName == owner {
					// println(rawCommand)
					params := strings.TrimPrefix(rawCommand, "manual(")
					params = strings.TrimSuffix(params, ")")
					println(params)
					orders := strings.Split(params, ",")
					for _, o := range orders {
						oi, err := strconv.ParseInt(o, 10, 64)
						if err != nil {
							log.Println(err)
							session.SendText("ERROR!What are U fucking typing! plz re-type it!", session.Bot.UserName, msg.FromUserName)
						}
						orderList = append(orderList, int(oi))
					}
					orderListStr := showRecipe(orderList)
					session.SendText(orderListStr, session.Bot.UserName, msg.FromUserName)
				} else {
					session.SendText("WTF!Dont kidding me!", session.Bot.UserName, msg.FromUserName)
				}
			} else if strings.Contains(rawCommand, "random") {
				if taking && msg.FromUserName == owner {
					params := strings.TrimSpace(strings.TrimPrefix(rawCommand, "random"))
					count, err := strconv.ParseInt(params, 10, 64)
					if err != nil {
						log.Println(err)
						session.SendText("ERROR!What are U fucking typing! plz re-type it!", session.Bot.UserName, msg.FromUserName)
					}
					orderListsName := getRecipe(int(count))
					if len(orderListsName) > 0 {
						session.SendText(showRecipeByName(orderListsName)+"\n type /send to send this order to BAOCAN", session.Bot.UserName, msg.FromUserName)
					}
					orderList = convertRecipe(orderListsName)
				} else {
					session.SendText("WTF!Dont kidding me!", session.Bot.UserName, msg.FromUserName)
				}
			} else if strings.Contains(rawCommand, "star") {
				if lunching && msg.FromUserName == owner {
					star := strings.TrimSpace(strings.TrimPrefix(rawCommand, "star"))
					starInt, err := strconv.ParseInt(star, 10, 64)
					if err != nil {
						log.Println(err)
						session.SendText("ERROR!What are U fucking typing! plz re-type it!", session.Bot.UserName, msg.FromUserName)
					}
					err = starRecipe(orderList, int(starInt))
					if err != nil {
						log.Println(err)
						session.SendText("ERROR!What are U fucking typing! plz re-type it!", session.Bot.UserName, msg.FromUserName)
					}
					session.SendText("The recipe is "+star+" ✨,it will work next time.type /end to end.", session.Bot.UserName, msg.FromUserName)
					owner = ""
				} else {
					session.SendText("WTF!Dont kidding me!", session.Bot.UserName, msg.FromUserName)
				}
			} else {
				session.SendText("WTF!", session.Bot.UserName, msg.FromUserName)
			}
		}
	}
}
