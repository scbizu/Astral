// Package tl provides TeamLiquid API wrappers
package tl

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	cache "github.com/patrickmn/go-cache"
	"github.com/robfig/cron"
	"github.com/scbizu/Astral/talker"
	"github.com/sirupsen/logrus"
)

const (
	timelineCacheKey = "timelines"
)

type timelines []int64

func (t timelines) getTheLastestTimeline() int64 {
	if len(t) == 0 {
		return 0
	}
	sort.SliceStable([]int64(t), func(i int, j int) bool {
		return []int64(t)[i] < []int64(t)[j]
	})
	return []int64(t)[0]
}

var (
	matchCache = cache.New(6*time.Hour, 12*time.Hour)
)

type mCron struct {
	ctx context.Context
	c   *cron.Cron
}

func NewCron() *mCron {
	return &mCron{
		ctx: context.TODO(),
		c:   cron.New(),
	}
}

type Fetcher struct {
	c     *mCron
	cache *cache.Cache
	Bot   *tgbotapi.BotAPI
}

func NewFetcher(bot *tgbotapi.BotAPI) *Fetcher {
	return &Fetcher{
		c:     new(mCron),
		cache: matchCache,
		Bot:   bot,
	}
}

func (f *Fetcher) Do() error {
	f.c = NewCron()
	f.c.c.AddFunc("@every 1m", func() {
		if f.cache.ItemCount() > 0 {
			now := time.Now()
			timeLines, ok := f.cache.Get(timelineCacheKey)
			if !ok {
				logrus.Errorf("get timeline cache failed: %s", "no cache key")
				return
			}
			cn, err := time.LoadLocation("Asia/Shanghai")
			if err != nil {
				logrus.Errorf("tl load location failed: %s", err.Error())
				return
			}
			t := new(timelines)
			if err = json.Unmarshal([]byte(timeLines.(string)), t); err != nil {
				return
			}
			if now.In(cn).Unix() < t.getTheLastestTimeline() {
				logrus.Infof("now is %d, the lasted match is at %d, no need to refresh cache.",
					now.In(cn).Unix(), t.getTheLastestTimeline())
				return
			}
		}

		logrus.Infof("warming TL cache...")
		if err := f.refreshCache(); err != nil {
			logrus.Errorf("refresh cache failed: %s", err.Error())
		}
	})
	f.c.c.Start()
	return nil
}

func (f *Fetcher) refreshCache() error {
	f.expireAllMatches()
	p, err := NewMatchParser()
	if err != nil {
		return err
	}
	timelines, err := p.GetTimelines()
	if err != nil {
		return err
	}
	tlJSON, err := json.Marshal(timelines)
	if err != nil {
		return err
	}
	f.cache.Set(timelineCacheKey, tlJSON, -1)
	matches, err := p.GetTimeMatches()
	if err != nil {
		return err
	}

	go f.pushMSG(timelines, matches)

	for t, m := range matches {
		f.cache.Set(strconv.FormatInt(t, 10), m, -1)
	}
	return nil
}

func (f *Fetcher) expireAllMatches() {
	f.cache.Flush()
}

func (f *Fetcher) pushMSG(tls []int64, matches map[int64][]Match) {
	sort.SliceStable(tls, func(i, j int) bool {
		return tls[i] < tls[j]
	})

	var sortedMatches []string
	for _, tl := range tls {
		ms, ok := matches[tl]
		if !ok {
			continue
		}
		for _, m := range ms {
			sortedMatches = append(sortedMatches, m.GetMDMatchInfo())
		}
	}
	f.pushWithLimit(sortedMatches, 5)
}

func (f *Fetcher) pushWithLimit(matches []string, limit int) {
	splitMatches := split(matches, limit)
	for _, sm := range splitMatches {
		matchPush := talker.NewMatchPush(sm)
		f.Bot.Send(matchPush.GetPushMessage())
	}
}

func split(buf []string, lim int) [][]string {
	var chunk []string
	chunks := make([][]string, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
