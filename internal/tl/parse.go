package tl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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
	matchesURL            = `https://liquipedia.net/starcraft2/api.php?action=parse&format=json&page=Liquipedia:Upcoming_and_ongoing_matches`
	timeFmt               = `January 2, 2006 - 15:04 UTC`
	maxCountDown          = time.Hour
	matchDetailFromIndex  = 19
	matchDetailEndIndex   = 21
	matchDetailPriceIndex = 17
	matchDetailOrganizer  = 3
)

type Timeline struct {
	T         int64 `json:"t"`
	IsOnGoing bool  `json:"isOnGoing"`
}

type MatchParser struct {
	rawHTML string
	revID   int
}

type Match struct {
	isOnGoing        bool
	vs               string
	timeCountingDown string
	series           string
	stream           []string
}

func (m Match) GetMDMatchInfo() string {
	if m.isOnGoing {
		return fmt.Sprintf(" 【🐔 比赛对阵】 %s \n 【🏆 所属杯赛】 %s \n 【📺 比赛直播】 %s", m.vs, m.series, strings.Join(m.stream, "/"))
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

	revID, rawHTML, err := newParseRespFromReader(r)
	if err != nil {
		return MatchParser{}, nil
	}
	return MatchParser{
		revID:   revID,
		rawHTML: rawHTML,
	}, nil
}

func newParseRespFromReader(r io.Reader) (int, string, error) {

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, "", err
	}

	tlMatches := new(TLMatchPage)

	if err := json.Unmarshal(body, tlMatches); err != nil {
		return 0, "", err
	}

	return tlMatches.Parse.Revid, tlMatches.Parse.Text.RawHTML, nil
}

func (mp MatchParser) GetRevID() int {
	return mp.revID
}

func (mp MatchParser) GetTimelines() ([]Timeline, error) {
	matches, err := mp.GetTimeMatches()
	if err != nil {
		return nil, err
	}

	ts := []Timeline{}
	for t, matches := range matches {
		var isOnGoing bool
		for _, m := range matches {
			if m.isOnGoing {
				isOnGoing = m.isOnGoing
				break
			}
		}
		ts = append(ts, Timeline{
			T:         t,
			IsOnGoing: isOnGoing,
		})
	}
	return ts, nil
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
			versus := s.Find(`.versus`).Text()
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
			if int64(countDown) <= 0 {
				var streams []string
				tournament := s.Find(`.match-filler > div`).Text()
				if tournament == "" {
					tournament = "未知"
				} else {
					detail := s.Find(`.match-filler > div > div > a`)
					if detail.Length() == 0 {
						logrus.Warn("match parser: match detail node not found")
					} else {
						detailURL, ok := detail.Attr("href")
						if !ok {
							matches[t.In(cn).Unix()] = append(matches[t.In(cn).Unix()], Match{
								isOnGoing: true,
								vs:        fmt.Sprintf("%s vs %s (%s)", trimText(lp), trimText(rp), versus),
								series:    strings.TrimSpace(tournament),
								stream:    []string{"直播源解析失败"},
							})
							return
						}
						u, err := url.Parse("https://liquipedia.net" + detailURL)
						if err != nil {
							logrus.Warnf("match parser: %q", err)
							matches[t.In(cn).Unix()] = append(matches[t.In(cn).Unix()], Match{
								isOnGoing: true,
								vs:        fmt.Sprintf("%s vs %s (%s)", trimText(lp), trimText(rp), versus),
								series:    strings.TrimSpace(tournament),
								stream:    []string{"直播源解析失败"},
							})
							return
						}
						md, err := getMatchDetail(u)
						if err != nil {
							logrus.Warnf("fetch match detail: %q", err)
							matches[t.In(cn).Unix()] = append(matches[t.In(cn).Unix()], Match{
								isOnGoing: true,
								vs:        fmt.Sprintf("%s vs %s (%s)", trimText(lp), trimText(rp), versus),
								series:    strings.TrimSpace(tournament),
								stream:    []string{"直播源解析失败"},
							})
							return
						}
						for _, s := range md.GetStreams() {
							if strings.TrimSpace(s.FmtToMarkdown()) == "" {
								continue
							}
							streams = append(streams, s.FmtToMarkdown())
						}
					}
				}
				logrus.Debugf("streams: %#v", streams)
				if len(streams) == 0 {
					streams = append(streams, "直播源解析失败")
				}
				vs := strings.Replace(versus, "\n", "", -1)
				if strings.Contains(vs, "vs") {
					vs = ""
				}
				matches[t.In(cn).Unix()] = append(matches[t.In(cn).Unix()], Match{
					isOnGoing: true,
					vs: fmt.Sprintf("%s vs %s (%s)",
						trimText(lp),
						trimText(rp),
						vs,
					),
					series: strings.TrimSpace(tournament),
					stream: streams,
				})
			}

			if 0 < int64(countDown) && int64(countDown) < int64(maxCountDown) {
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

type MatchDetail struct {
	players   []string
	from      time.Time
	to        time.Time
	prize     string
	organizer string
	streams   []Stream
}

func getMatchDetail(u *url.URL) (MatchDetail, error) {
	mURL := u.String()
	logrus.Debugf("match full url: %s", u.String())
	d, err := goquery.NewDocument(mURL)
	if err != nil {
		return MatchDetail{}, err
	}
	var start, end time.Time
	var prize string
	var orger string
	d.Find(`.infobox-cell-2`).Each(func(index int, s *goquery.Selection) {
		if index == matchDetailFromIndex {
			st, err := time.Parse(`2006-01-02`, s.Text())
			if err != nil {
				logrus.Warnf("match detail: %q", err.Error())
				return
			}
			start = st
		}

		if index == matchDetailEndIndex {
			ed, err := time.Parse(`2006-01-02`, s.Text())
			if err != nil {
				logrus.Warnf("match detail: %q", err.Error())
				return
			}
			end = ed
		}

		if index == matchDetailPriceIndex {
			prize = s.Text()
		}

		if index == matchDetailOrganizer {
			orger = s.Text()
		}
	})

	var streams []Stream
	d.Find(`#Streams`).Parent().Next().Find(`li > a.external.text`).
		Each(func(index int, s *goquery.Selection) {
			logrus.Debugf("match streaming: %s", s.Text())
			stream := Stream{}
			h, ok := s.Attr(`href`)
			if ok {
				stream.caster = s.Text()
				stream.streammingURL = h
			} else {
				stream.caster = "unknown"
			}
			streams = append(streams, stream)
		})

	return MatchDetail{
		streams:   streams,
		from:      start,
		to:        end,
		prize:     prize,
		organizer: orger,
		// TODO: parse players
		players: []string{},
	}, nil
}

func (md MatchDetail) GetStreams() []Stream {
	return md.streams
}

type Stream struct {
	caster        string
	streammingURL string
}

func (s Stream) FmtToMarkdown() string {
	return fmt.Sprintf("[%s](%s)", s.caster, s.streammingURL)
}
