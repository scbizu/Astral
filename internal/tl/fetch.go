// Package tl provides TeamLiquid API wrappers
package tl

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/robfig/cron"
	"github.com/scbizu/Astral/internal/mcache"
	"github.com/sirupsen/logrus"
)

const (
	timelineCacheKey = "timelines"
	version          = "revID"
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
	c             *mCron
	cache         *cache.Cache
	dsts          []IRC
	pushedMatches []Match
	sync.Mutex
}

func NewFetcher(s ...IRC) *Fetcher {
	return &Fetcher{
		c:     new(mCron),
		cache: matchCache,
		dsts:  s,
	}
}

func (f *Fetcher) Register() {
	// register match message cache
	mcache.EnableMessageCache()
}

func (f *Fetcher) Do() error {

	f.Register()

	f.c = NewCron()
	if err := f.c.c.AddFunc("@every 5m", func() {
		if err := f.refreshCache(); err != nil {
			logrus.Errorf("refresh cache failed: %s", err.Error())
		}
	}); err != nil {
		return fmt.Errorf("cron: %q", err)
	}
	f.c.c.Run()
	return nil
}

func (f *Fetcher) refreshCache() error {
	p, err := NewMatchParser()
	if err != nil {
		return err
	}
	cacheRevID, ok := f.cache.Get(version)
	if ok && cacheRevID == p.GetRevID() {
		return nil
	}
	logrus.Infof("LastRevID: %d, CurrentRevID: %d", cacheRevID, p.GetRevID())
	defer func() {
		f.cache.Set(version, p.GetRevID(), -1)
	}()
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

	for t, ms := range matches {
		// event is still alive
		if _, ok := f.cache.Get(strconv.FormatInt(t, 10)); ok {
			continue
		}
		f.cache.Set(strconv.FormatInt(t, 10), ms, -1)
	}
	return nil
}

func (f *Fetcher) pushMSG(tls []Timeline, matches map[int64][]Match) {
	sort.SliceStable(tls, func(i, j int) bool {
		return tls[i].T < tls[j].T
	})

	var sortedMatches []Match
	for _, tl := range tls {
		// matches must be the superset of the tls
		ms, ok := matches[tl.T]
		if !ok {
			continue
		}
		sortedMatches = append(sortedMatches, ms...)
	}
	f.pushWithLimit(sortedMatches, 5)
}

func (f *Fetcher) pushWithLimit(matches []Match, limit int) {

	var matchStr []string

	for _, m := range matches {
		matchStr = append(matchStr, m.GetMDMatchInfo())
	}

	splitMatchesStr := split(matchStr, limit)
	if len(splitMatchesStr) == 0 {
		return
	}

	splitMatches := splitMatch(matches, limit)

	f.pushedMatches = append(f.pushedMatches, matches...)

	// use n goroutines to send message
	for _, dst := range f.dsts {
		go func(dst Sender) {
			var idx int
		SEND:
			msg := dst.ResolveMessage(splitMatchesStr[idx])
			if err := dst.Send(msg, CacheMessageFilter{}); err != nil {
				logrus.Errorf("sender: %s", err.Error())
				return
			}
			f.Lock()
			if err := mcache.AddMessage(msg, mcache.MD5{}); err != nil {
				logrus.Errorf("mcache: set cache: %q", err)
				// ignore and fallthrough
			}
			f.Unlock()
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
		chunks = append(chunks, buf[:])
	}
	return chunks
}

func splitMatch(buf []Match, lim int) [][]Match {
	var chunk []Match
	chunks := make([][]Match, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}
