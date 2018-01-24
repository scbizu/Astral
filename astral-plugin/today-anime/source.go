package anime

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//SrcObj defines bangumi obj
type SrcObj struct {
	Src         string
	BangumiName string
	Link        *url.URL
	Pubed       bool
}

const (
	//BilibiliGC B站国创
	BilibiliGC = "https://bangumi.bilibili.com/guochuang/timeline"
	//BilibiliJP B站日漫
	BilibiliJP = "https://bangumi.bilibili.com/anime/timeline"
	//Dilidili D站动漫
	Dilidili = "http://www.dilidili.wang"
)

//FormatLinkInMarkdownPreview formats srcobj to Markdown view
func (s *SrcObj) FormatLinkInMarkdownPreview() string {
	name := fmt.Sprintf("[%s From %s]", s.BangumiName, s.Src)
	linkstr := fmt.Sprintf("(%s)", s.Link.String())
	return fmt.Sprintf("%s%s", name, linkstr)
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
	resURL := fmt.Sprintf("%s%s", rawlink, src)
	resLink, err = url.Parse(resURL)
	return
}

func scrapeBilibiliTimeline(src *url.URL) ([]*SrcObj, error) {
	doc, err := goquery.NewDocument(src.String())
	if err != nil {
		return nil, err
	}
	var objs []*SrcObj
	log.Print("fetching...")
	content, _ := doc.Find("body").Html()
	log.Print(content)
	doc.Find(".day-wrap current").Each(func(index int, s *goquery.Selection) {
		log.Print(index)
		obj := new(SrcObj)
		obj.Src = "bilibili"
		link, _ := s.Find(".tl-body a").Eq(1).Attr("href")
		obj.Link, err = formatLink(link)
		if err != nil {
			log.Printf("format bilibili url error:%s", err.Error())
		}
		obj.BangumiName, _ = s.Find(".tl-body a").Eq(1).Attr("title")
		pubstr := s.Find(".tl-body a").Eq(1).Find(".published").Text()
		if pubstr == "" {
			obj.Pubed = false
		} else {
			obj.Pubed = true
		}
		objs = append(objs, obj)
	})
	return objs, nil
}

func scrapeDilidiliTimeLine(src *url.URL) ([]*SrcObj, error) {
	doc, err := goquery.NewDocument(src.String())
	if err != nil {
		return nil, err
	}
	var objs []*SrcObj

	doc.Find(".container-row-1 update .two-auto ul").Each(func(index int, s *goquery.Selection) {
		v, _ := s.Attr("style")
		if v == "display: block" {
			s.Find(".tooltip tooltipstered").Each(func(cindex int, cs *goquery.Selection) {
				ele := cs.Find(".update-content h4 a")
				obj := new(SrcObj)
				obj.Src = "dilidili"
				obj.BangumiName = ele.Text()
				link, _ := ele.Attr("href")
				var err error
				obj.Link, err = formatNotAbsoluteLink(link, Dilidili)
				if err != nil {
					log.Printf("format dilidili url error:%s", err.Error())
				}
				obj.Pubed = false
				objs = append(objs, obj)
			})
		}
	})
	return objs, nil
}
