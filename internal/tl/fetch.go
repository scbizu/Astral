// Package tl provides TeamLiquid API wrappers
package tl

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

const (
	timelineCacheKey = "timelines"
)

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
	dsts  []Sender
}

func NewFetcher(s ...Sender) *Fetcher {
	return &Fetcher{
		c:     new(mCron),
		cache: matchCache,
		dsts:  s,
	}
}

func (f *Fetcher) Do() error {
	f.c = NewCron()
	f.c.c.AddFunc("@every 1m", func() {
		// do some condition tricks to save bandwidth
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

			t := make([]Timeline, 0)
			if err = json.Unmarshal(timeLines.([]byte), &t); err != nil {
				logrus.Errorf("unmarshal cache failed: %s", err.Error())
				return
			}

			if len(t) == 0 {
				logrus.Warn("get 0 timelines")
				return
			}

			if now.In(cn).Unix() < getTheLastestTimeline(t) {
				logrus.Infof("now is %d, the lastest match is at %d, no need to refresh cache.",
					now.In(cn).Unix(), getTheLastestTimeline(t))
				return
			}
		}

		logrus.Infof("warming TL cache...")
		if err := f.refreshCache(); err != nil {
			logrus.Errorf("refresh cache failed: %s", err.Error())
		}
	})
	f.c.c.Run()
	return nil
}

func (f *Fetcher) refreshCache() error {
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

	go f.pushMSG(timelines, f.expireCache(timelines, matches))

	for t, m := range matches {
		f.cache.Set(strconv.FormatInt(t, 10), m, -1)
	}
	return nil
}

// expireCache reuse the ongoing match info
// and delete the out-of-date match info
// Due to the reuseable cache, from now on , we should manage our cache carefully TAT
func (f *Fetcher) expireCache(tls []Timeline, matches map[int64][]Match) map[int64][]Match {
	// reuse cache:
	// TL will reset Streaming caster URL after the match is going.
	// We should keep the opening match info until it is closed.
	for t := range matches {
		cachedMatches, ok := f.cache.Get(strconv.FormatInt(t, 10))
		if ok {
			matches[t] = cachedMatches.([]Match)
		}
	}

	if len(tls) == 0 {
		return matches
	}

	sort.SliceStable(tls, func(i, j int) bool {
		return tls[i].T < tls[j].T
	})

	// expire cache : T is less than the index 0 (the lowest one)
	for t := range matches {
		if t >= tls[0].T {
			continue
		}
		delete(matches, t)
		f.cache.Delete(strconv.FormatInt(t, 10))
	}
	return matches
}

func (f *Fetcher) pushMSG(tls []Timeline, matches map[int64][]Match) {
	sort.SliceStable(tls, func(i, j int) bool {
		return tls[i].T < tls[j].T
	})

	var sortedMatches []string
	for _, tl := range tls {
		// matches must be the superset of the tls
		ms, ok := matches[tl.T]
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
	if len(splitMatches) == 0 {
		return
	}
	// use n goroutines to send message
	for _, dst := range f.dsts {
		go func(dst Sender) {
			var idx int
		SEND:
			msg := dst.ResolveMessage(splitMatches[idx])
			if err := dst.Send(msg); err != nil {
				logrus.Errorf("sender: %s", err.Error())
			}
			idx++
			if idx < len(splitMatches)-1 {
				goto SEND
			}
		}(dst)
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

func getTheLastestTimeline(tls []Timeline) int64 {
	sort.SliceStable(tls, func(i, j int) bool {
		return tls[i].T < tls[j].T
	})

	for _, tl := range tls {
		if !tl.IsOnGoing {
			return tl.T
		}
	}

	return 0
}
