package anime

import (
	"net/url"
	"testing"
)

func TestScrapeBilibiliTimeline(t *testing.T) {
	srcURL, err := url.Parse(BilibiliGC)
	if err != nil {
		t.Error(err)
		return
	}
	infos, err := scrapeBilibiliTimeline(srcURL)
	if err != nil {
		t.Error(err)
		return
	}
	for _, info := range infos {
		t.Logf("AnimeName:%s", info.BangumiName)
	}
}
