package tl

import "github.com/scbizu/Astral/internal/mcache"

type CacheMessageFilter struct{}

func (CacheMessageFilter) F(raw string) string {
	if mcache.IsMessageSet(raw, mcache.MD5{}) {
		return ""
	}
	return raw
}
