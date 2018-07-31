package anime

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//TimeLine defines bilibili's timeline
type TimeLine struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  []struct {
		Date      string `json:"date"`
		DateTs    int    `json:"date_ts"`
		DayOfWeek int    `json:"day_of_week"`
		IsToday   int    `json:"is_today"`
		Seasons   []struct {
			Cover        string `json:"cover"`
			Delay        int    `json:"delay"`
			EpID         int    `json:"ep_id"`
			Favorites    int    `json:"favorites"`
			Follow       int    `json:"follow"`
			IsPublished  int    `json:"is_published"`
			PubIndex     string `json:"pub_index"`
			PubTime      string `json:"pub_time"`
			PubTs        int    `json:"pub_ts"`
			SeasonID     int    `json:"season_id"`
			SeasonStatus int    `json:"season_status"`
			SquareCover  string `json:"square_cover"`
			Title        string `json:"title"`
			Badge        string `json:"badge,omitempty"`
		} `json:"seasons"`
	} `json:"result"`
}

//TimeLineCN defines 国创 timeline
type TimeLineCN struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  []struct {
		Date      string `json:"date"`
		DateTs    int    `json:"date_ts"`
		DayOfWeek int    `json:"day_of_week"`
		IsToday   int    `json:"is_today"`
		Seasons   []struct {
			Cover        string `json:"cover"`
			Delay        int    `json:"delay"`
			EpID         int    `json:"ep_id"`
			Favorites    int    `json:"favorites"`
			Follow       int    `json:"follow"`
			IsPublished  int    `json:"is_published"`
			PubIndex     string `json:"pub_index"`
			PubTime      string `json:"pub_time"`
			PubTs        int    `json:"pub_ts"`
			SeasonID     int    `json:"season_id"`
			SeasonStatus int    `json:"season_status"`
			SquareCover  string `json:"square_cover"`
			Title        string `json:"title"`
		} `json:"seasons"`
	} `json:"result"`
}

//SrcObj defines bangumi obj
type SrcObj struct {
	Src         string
	BangumiName string
	Link        *url.URL
	Pubed       bool
}

const (
	//BilibiliGC B站国创
	BilibiliGC = "https://bangumi.bilibili.com/web_api/timeline_cn"
	//BilibiliJP B站日漫
	BilibiliJP = "https://bangumi.bilibili.com/web_api/timeline_global"
	//Dilidili D站动漫
	Dilidili = "http://www.dilidili.wang"
)

//FormatLinkInMarkdownPreview formats srcobj to Markdown view
func (s *SrcObj) FormatLinkInMarkdownPreview() string {
	name := fmt.Sprintf("[%s From %s]", s.BangumiName, s.Src)
	linkstr := fmt.Sprintf("(%s)", s.Link.String())
	if s.Pubed {
		return fmt.Sprintf("%s%s", name, linkstr)
	}
	return fmt.Sprintf("%s From %s(未更新)", s.BangumiName, s.Src)
}

//GetAllAnimes gets all animes from all src defined.
func GetAllAnimes() (objs []*SrcObj, err error) {
	objs, err = GetAnimeFromBGC()
	if err != nil {
		return
	}
	bbjp, err := GetAnimeFromBJP()
	if err != nil {
		return nil, err
	}
	objs = append(objs, bbjp...)
	d, err := GetAnimeFromD()
	if err != nil {
		return nil, err
	}
	objs = append(objs, d...)
	return
}

// GetAnimeFromB get anime from bilibili
func GetAnimeFromB() (objs []*SrcObj, err error) {
	objs, err = GetAnimeFromBGC()
	if err != nil {
		return
	}
	bbjp, err := GetAnimeFromBJP()
	if err != nil {
		return nil, err
	}
	objs = append(objs, bbjp...)
	return
}

//GetAnimeFromBGC ....
func GetAnimeFromBGC() ([]*SrcObj, error) {
	bgcSrc, err := url.Parse(BilibiliGC)
	if err != nil {
		return nil, err
	}
	objs, err := scrapeBilibiliTimeline(bgcSrc)
	if err != nil {
		return nil, err
	}

	return objs, nil
}

//GetAnimeFromBJP ...
func GetAnimeFromBJP() ([]*SrcObj, error) {
	bjpSrc, err := url.Parse(BilibiliJP)
	if err != nil {
		return nil, err
	}
	objs, err := scrapeBilibiliTimeline(bjpSrc)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

//GetAnimeFromD ...
func GetAnimeFromD() ([]*SrcObj, error) {
	diliURL, err := url.Parse(Dilidili)
	if err != nil {
		return nil, err
	}
	objs, err := scrapeDilidiliTimeLine(diliURL)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

func formatLink(rawLink string) (resLink *url.URL, err error) {
	var resURL string
	if strings.HasPrefix(rawLink, "//") {
		resURL = fmt.Sprintf("https:%s", rawLink)
	} else {
		resURL = rawLink
	}
	resLink, err = url.Parse(resURL)
	if err != nil {
		return
	}
	return
}

func formatNotAbsoluteLink(rawlink string, src string) (resLink *url.URL, err error) {
	resURL := fmt.Sprintf("%s%s", src, rawlink)
	resLink, err = url.Parse(resURL)
	return
}

func scrapeBilibiliTimeline(src *url.URL) ([]*SrcObj, error) {
	req, err := http.Get(src.String())
	if err != nil {
		return nil, err
	}
	var objs []*SrcObj
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	tl := new(TimeLine)
	err = json.Unmarshal(body, tl)
	if err != nil {
		return nil, err
	}
	for _, r := range tl.Result {
		if r.IsToday == 0 {
			continue
		}
		for _, s := range r.Seasons {
			if s.Delay == 1 {
				continue
			}
			obj := new(SrcObj)
			obj.BangumiName = s.Title
			obj.Link, err = formatNotAbsoluteLink(strconv.Itoa(s.EpID),
				"https://www.bilibili.com/bangumi/play/ep")
			if err != nil {
				return nil, err
			}
			obj.Pubed = s.IsPublished > 0
			obj.Src = "bilibili"
			objs = append(objs, obj)
		}
	}

	return objs, nil
}

func scrapeDilidiliTimeLine(src *url.URL) ([]*SrcObj, error) {
	doc, err := goquery.NewDocument(src.String())
	if err != nil {
		return nil, err
	}
	var objs []*SrcObj
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}

	today := convert2CNWeekDay(int(time.Now().In(location).Weekday()), 6)
	// log.Println(doc.Find(".container-row-1").Find(".two-auto").Find("ul").Find(''))

	doc.Find(".change").Eq(1).Find(".sldr").Find(".wrp > li").Each(func(index int, s *goquery.Selection) {
		if index == today {
			s.Find(".list > li").Each(func(cindex int, cs *goquery.Selection) {
				ele := cs.Find("a")
				obj := new(SrcObj)
				obj.Src = "dilidili"
				obj.BangumiName = ele.Text()
				link, _ := ele.Attr("href")
				var err error
				obj.Link, err = formatNotAbsoluteLink(link, Dilidili)
				if err != nil {
					log.Printf("format dilidili url error:%s", err.Error())
				}
				obj.Pubed = true
				objs = append(objs, obj)
			})
		}
	})
	return objs, nil
}

// weekday internationalWeekday    cnWeekday
// Sun     0      6
// Mon     1      0
//`Tue     2      1
// Wed     3      2
// Thu     4      3
// Fri     5      4
// Sat     6      5
func convert2CNWeekDay(internationWeekday int, offsetDay int) (cnWeekday int) {
	return (internationWeekday + offsetDay) % 7
}
