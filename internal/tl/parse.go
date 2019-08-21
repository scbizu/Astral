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
	"github.com/scylladb/go-set/strset"
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
	maxCountDown          = 20 * time.Minute
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
	vs               Versus
	timeCountingDown string
	series           string
	stream           []string
	detailURL        *url.URL
}

func (m Match) GetVS() string {
	return m.vs.f()
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
			vs := strings.Replace(versus, "\n", "", -1)
			if strings.Contains(vs, "vs") {
				vs = ""
			}
			score := strings.Split(versus, ":")
			var s1, s2 string
			if len(score) < 2 {
				s1 = "0"
				s2 = "0"
			} else {
				s1 = score[0]
				s2 = score[1]
			}
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
				var u *url.URL
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
								vs: Versus{
									P1:      trimText(lp),
									P2:      trimText(rp),
									P1Score: s1,
									P2Score: s2,
								},
								series: strings.TrimSpace(tournament),
								stream: []string{"直播源解析失败"},
							})
							return
						}
						var err error
						u, err = url.Parse("https://liquipedia.net" + detailURL)
						if err != nil {
							logrus.Warnf("match parser: %q", err)
							matches[t.In(cn).Unix()] = append(matches[t.In(cn).Unix()], Match{
								isOnGoing: true,
								vs: Versus{
									P1:      trimText(lp),
									P2:      trimText(rp),
									P1Score: s1,
									P2Score: s2,
								},
								series: strings.TrimSpace(tournament),
								stream: []string{"直播源解析失败"},
							})
							return
						}
						md, err := getMatchDetail(u)
						if err != nil {
							logrus.Warnf("fetch match detail: %q", err)
							matches[t.In(cn).Unix()] = append(matches[t.In(cn).Unix()], Match{
								isOnGoing: true,
								vs: Versus{
									P1:      trimText(lp),
									P2:      trimText(rp),
									P1Score: s1,
									P2Score: s2,
								},
								series: strings.TrimSpace(tournament),
								stream: []string{"直播源解析失败"},
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
				matches[t.In(cn).Unix()] = append(matches[t.In(cn).Unix()], Match{
					isOnGoing: true,
					vs: Versus{
						P1:      trimText(lp),
						P2:      trimText(rp),
						P1Score: s1,
						P2Score: s2,
					},
					series:    strings.TrimSpace(tournament),
					stream:    streams,
					detailURL: u,
				})
			}

			if 0 < int64(countDown) && int64(countDown) < int64(maxCountDown) {
				tournament := s.Find(`.matchticker-tournament-wrapper`).Text()
				if tournament == "" {
					tournament = "未知"
				}
				matches[t.In(cn).Unix()] = append(matches[t.In(cn).Unix()], Match{
					isOnGoing: false,
					vs: Versus{
						P1:      trimText(lp),
						P2:      trimText(rp),
						P1Score: s1,
						P2Score: s2,
					},
					timeCountingDown: countDown.String(),
					series:           strings.TrimSpace(tournament),
				})
			}
		})
	return matches, nil
}

func GetHTMLMatchesResp() (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", matchesURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "AstralBot(https://github.com/scbizu/Astral)")
	resp, err := http.DefaultClient.Do(req)
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

type Versus struct {
	P1      string
	P2      string
	P1Score string
	P2Score string
}

func (v Versus) f() string {
	return fmt.Sprintf("%s vs %s (%s:%s)", v.P1, v.P2, v.P1Score, v.P2Score)
}

func GetFinalMatchRes(u *url.URL, p1, p2 string) (Versus, error) {
	d, err := goquery.NewDocument(u.String())
	if err != nil {
		return Versus{}, err
	}
	vs := Versus{}
	std := strset.New()
	std.Add(p1, p2)
	d.Find(`.matchlistslot`).Each(func(index int, s *goquery.Selection) {
		if trimText(s.Text()) == p1 || trimText(s.Text()) == p2 {
			row := s.Parent()
			gets := strset.New()
			playerScore := make(map[string]string)
			row.Find(`td`).Each(func(subIndex int, s *goquery.Selection) {
				switch subIndex {
				case 0, 2:
					gets.Add(trimText(s.Text()))
					// case 1,3
					playerScore[trimText(s.Text())] = trimText(s.Next().Text())
				}
			})
			if gets.IsEqual(std) {
				// find in Group Stage
				vs.P1 = gets.Pop()
				vs.P2 = gets.Pop()
				var ok bool
				vs.P1Score, ok = playerScore[vs.P1]
				if !ok {
					vs.P1Score = "0"
				}
				vs.P2Score, ok = playerScore[vs.P2]
				if !ok {
					vs.P2Score = "0"
				}
			}
		}
	})

	// Keep always getting the lastest versus info, it will ignore Group Stage versus information, but it is correct.
	// Playoffs
	d.Find(`.bracket-cell-r1`).Each(func(index int, s *goquery.Selection) {
		if index%2 == 0 {
			return
		}
		if strings.Contains(trimText(s.Text()), p1) || strings.Contains(trimText(s.Text()), p2) {
			row := s.Parent()
			gets := strset.New()
			playerScore := make(map[string]string)
			row.Find(`.bracket-cell-r1`).Each(func(subIndex int, sub *goquery.Selection) {
				score := trimText(sub.Find(`.bracket-score`).Text())
				// e.g: TIME0 => map{"TIME":"0"}
				player := strings.TrimSuffix(trimText(sub.Text()), score)
				gets.Add(player)
				playerScore[player] = score
			})

			if gets.IsEqual(std) {
				vs.P1 = gets.Pop()
				vs.P2 = gets.Pop()
				var ok bool
				vs.P1Score, ok = playerScore[vs.P1]
				if !ok {
					vs.P1Score = "0"
				}
				vs.P2Score, ok = playerScore[vs.P2]
				if !ok {
					vs.P2Score = "0"
				}
			}
		}
	})

	return vs, nil
}
