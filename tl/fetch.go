// Package tl provides TeamLiquid API wrappers
package tl

import (
	"context"
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
		logrus.Infof("warming TL cache...")
		if f.cache.ItemCount() > 0 {
			now := time.Now()
			timeLines, ok := f.cache.Get(timelineCacheKey)
			if !ok {
				return
			}
			cn, err := time.LoadLocation("Asia/Shanghai")
			if err != nil {
				logrus.Errorf("tl load location failed: %s", err.Error())
				return
			}
			timeLineInts, ok := timeLines.(timelines)
			if !ok {
				return
			}
			if now.In(cn).Unix() < timeLineInts.getTheLastestTimeline() {
				return
			}
		}
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
	f.cache.Set(timelineCacheKey, timelines, 6*time.Hour)
	matches, err := p.GetTimeMatches()
	if err != nil {
		return err
	}

	go f.pushMSG(timelines, matches)

	for t, m := range matches {
		f.cache.Set(strconv.FormatInt(t, 10), m, 6*time.Hour)
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
	matchPush := talker.NewMatchPush(sortedMatches)
	f.Bot.Send(matchPush.GetPushMessage())
}
