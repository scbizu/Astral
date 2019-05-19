package tl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type TLMatchPage struct {
	Parse struct {
		Title  string `json:"title"`
		Pageid int    `json:"pageid"`
		Revid  int    `json:"revid"`
		Text   struct {
			RawHTML string `json:"*"`
		} `json:"text"`
		Langlinks  []interface{} `json:"langlinks"`
		Categories []interface{} `json:"categories"`
		Links      []struct {
			Ns     int    `json:"ns"`
			Exists string `json:"exists,omitempty"`
			All    string `json:"*"`
		} `json:"links"`
		Templates     []interface{} `json:"templates"`
		Images        []string      `json:"images"`
		Externallinks []interface{} `json:"externallinks"`
		Sections      []struct {
			Toclevel   int    `json:"toclevel"`
			Level      string `json:"level"`
			Line       string `json:"line"`
			Number     string `json:"number"`
			Index      string `json:"index"`
			Fromtitle  string `json:"fromtitle"`
			Byteoffset int    `json:"byteoffset"`
			Anchor     string `json:"anchor"`
		} `json:"sections"`
		Parsewarnings []interface{} `json:"parsewarnings"`
		Displaytitle  string        `json:"displaytitle"`
		Iwlinks       []interface{} `json:"iwlinks"`
		Properties    []struct {
			Name string `json:"name"`
			All  string `json:"*"`
		} `json:"properties"`
	} `json:"parse"`
}

const (
	matchesURL   = `https://liquipedia.net/starcraft2/api.php?action=parse&format=json&page=Liquipedia:Upcoming_and_ongoing_matches`
	timeFmt      = `January 2, 2006 - 15:04 UTC`
	maxCountDown = time.Hour
)

type Timeline struct {
	T         int64 `json:"t"`
	IsOnGoing bool  `json:"isOnGoing"`
}

type MatchParser struct {
	rawHTML string
}

type Match struct {
	isOnGoing        bool
	vs               string
	timeCountingDown string
	series           string
}

func (m Match) GetMDMatchInfo() string {
	if m.isOnGoing {
		return fmt.Sprintf(" 【🐔 比赛对阵】 %s \n 【🏆 所属杯赛】 %s \n 【⏳ 比赛状态】 正在进行", m.vs, m.series)
	}
	return fmt.Sprintf(" 【🐔 比赛对阵】 %s \n 【🏆 所属杯赛】 %s \n 【⏳ 比赛状态】 倒计时 %s", m.vs, m.series, m.timeCountingDown)
}

func (m Match) GetJSONMatchInfo() (string, error) {
	matchesJSON, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(matchesJSON), nil
}

func NewMatchParser() (MatchParser, error) {
	r, err := GetHTMLMatchesResp()
	if err != nil {
		return MatchParser{}, err
	}
	defer r.Close()

	rawHTML, err := newParseRespFromReader(r)
	if err != nil {
		return MatchParser{}, nil
	}
	return MatchParser{
		rawHTML: rawHTML,
	}, nil
}

func newParseRespFromReader(r io.Reader) (string, error) {

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	tlMatches := new(TLMatchPage)

	if err := json.Unmarshal(body, tlMatches); err != nil {
		return "", err
	}

	return tlMatches.Parse.Text.RawHTML, nil
}

func (mp MatchParser) GetTimelines() ([]Timeline, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(mp.rawHTML))
	if err != nil {
		return nil, err
	}
	var timelines []Timeline
	doc.Find(`.timer-object-countdown-only`).Each(func(idx int, s *goquery.Selection) {
		timelineStd, err := time.Parse(timeFmt, s.Text())
		if err != nil {
			logrus.Errorf("parse failed: %s", err.Error())
			return
		}
		cn, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			logrus.Errorf("parse failed: %s", err.Error())
			return
		}

		countDown := time.Until(timelineStd.In(cn))
		if countDown > 0 {
			timelines = append(timelines, Timeline{
				IsOnGoing: false,
				T:         timelineStd.In(cn).Unix(),
			})
		} else if countDown < maxCountDown {
			timelines = append(timelines, Timeline{
				IsOnGoing: true,
				T:         timelineStd.In(cn).Unix(),
			})
		}
	})
	return timelines, nil
}

func (mp MatchParser) GetTimeMatches() (map[int64][]Match, error) {

	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(mp.rawHTML))
	if err != nil {
		return nil, err
	}
	matches := make(map[int64][]Match)
	doc.Find(`.infobox_matches_content`).
		Each(func(idx int, s *goquery.Selection) {
			lp := s.Find(`.team-left`).Text()
			rp := s.Find(`.team-right`).Text()
			t, err := time.Parse(timeFmt, s.Find(`.timer-object-countdown-only`).Text())
			if err != nil {
				logrus.Errorf("parse failed: %s", err.Error())
				return
			}
			cn, err := time.LoadLocation("Asia/Shanghai")
			if err != nil {
				logrus.Errorf("parse failed: %s", err.Error())
				return
			}
			countDown := time.Until(t.In(cn))
			if int64(countDown) < 0 {
				tournament := s.Find(`.match-filler > div`).Text()
				if tournament == "" {
					tournament = "未知"
				}
				matches[t.In(cn).Unix()] = append(matches[t.In(cn).Unix()], Match{
					isOnGoing:        true,
					vs:               fmt.Sprintf("%s vs %s", trimText(lp), trimText(rp)),
					timeCountingDown: "",
					series:           strings.TrimSpace(tournament),
				})
			} else if countDown < maxCountDown {
				tournament := s.Find(`.matchticker-tournament-wrapper`).Text()
				if tournament == "" {
					tournament = "未知"
				}
				matches[t.In(cn).Unix()] = append(matches[t.In(cn).Unix()], Match{
					isOnGoing:        false,
					vs:               fmt.Sprintf("%s vs %s", trimText(lp), trimText(rp)),
					timeCountingDown: countDown.String(),
					series:           strings.TrimSpace(tournament),
				})
			}
		})
	return matches, nil
}

func GetHTMLMatchesResp() (io.ReadCloser, error) {
	resp, err := http.Get(matchesURL)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func trimText(str string) string {
	return strings.TrimSpace(str)
}
