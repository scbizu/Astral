// Package tl provides TeamLiquid API wrappers
package tl

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/robfig/cron"
	"github.com/scylladb/go-set/strset"
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

type PMatch struct {
	msgID      string
	matchIndex int
	m          Match
}

type Fetcher struct {
	c             *mCron
	cache         *cache.Cache
	dsts          []Sender
	pushedMatches map[string]PMatch
	sync.Mutex
}

func NewFetcher(s ...Sender) *Fetcher {
	p := make(map[string]PMatch)
	return &Fetcher{
		c:             new(mCron),
		cache:         matchCache,
		dsts:          s,
		pushedMatches: p,
	}
}

func (f *Fetcher) Do() error {
	// run match GC
	go GetStashChan().Run()
	f.c = NewCron()
	f.c.c.AddFunc("@every 5m", func() {
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

	matches = f.reuseCache(timelines, matches)

	go f.pushMSG(timelines, matches)

	for t, m := range matches {
		f.cache.Set(strconv.FormatInt(t, 10), m, -1)
	}
	return nil
}

// reuseCache reuse the ongoing match info
// and delete the out-of-date match info
// Due to the reuseable cache, from now on , we should manage our cache carefully TAT
func (f *Fetcher) reuseCache(tls []Timeline, matches map[int64][]Match) map[int64][]Match {
	// reuse cache:
	// * the same versus information
	for t := range matches {
		var ms []Match
		// build caches versus
		vss := strset.New()
		cachedMatches, ok := f.cache.Get(strconv.FormatInt(t, 10))
		if ok {
			ms, mOK := cachedMatches.([]Match)
			if !mOK {
				continue
			}
			for _, m := range ms {
				vss.Add(m.GetVS())
			}
		}
		// filter matches
		for _, m := range matches[t] {
			if vss.Has(m.GetVS()) {
				continue
			}
			ms = append(ms, m)
		}
		matches[t] = ms
	}

	if len(tls) == 0 {
		return matches
	}

	sort.SliceStable(tls, func(i, j int) bool {
		return tls[i].T < tls[j].T
	})

	// expire cache : T is less than the index 0 (the lowest one)
	for t, ms := range matches {
		if t >= tls[0].T {
			continue
		}
		for _, m := range ms {
			pm, ok := f.pushedMatches[m.GetMDMatchInfo()]
			if !ok {
				continue
			}
			GetStashChan().Put(pm)
			// terminated
			f.Lock()
			delete(f.pushedMatches, m.GetMDMatchInfo())
			f.Unlock()
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

	var sortedMatches []Match
	for _, tl := range tls {
		// matches must be the superset of the tls
		ms, ok := matches[tl.T]
		if !ok {
			continue
		}
		for _, m := range ms {
			sortedMatches = append(sortedMatches, m)
		}
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

	// use n goroutines to send message
	for _, dst := range f.dsts {
		go func(dst Sender) {
			var idx int
		SEND:
			msg := dst.ResolveMessage(splitMatchesStr[idx])
			if err := dst.Send(msg); err != nil {
				logrus.Errorf("sender: %s", err.Error())
			}
			f.Lock()
			for i, m := range splitMatches[idx] {
				// TODO: msgid
				msg := ""
				f.pushedMatches[m.GetMDMatchInfo()] = PMatch{msgID: msg, matchIndex: i, m: m}
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
		chunks = append(chunks, buf[:len(buf)])
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
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
