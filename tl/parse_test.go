package tl

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMatchParser_GetTimelines(t *testing.T) {
	// r, err := ioutil.ReadFile("test.json")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// resp, err := newParseRespFromReader(bytes.NewBuffer(r))
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	mp, err := NewMatchParser()
	if err != nil {
		t.Error(err)
	}
	tls, err := mp.GetTimelines()
	if err != nil {
		t.FailNow()
	}
	if len(tls) > 0 {
		logrus.Infof("testing timelines out : %v", tls)
	}
}

func TestMatchParser_GetTimeMatches(t *testing.T) {
	r, err := ioutil.ReadFile("test.json")
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := newParseRespFromReader(bytes.NewBuffer(r))
	if err != nil {
		t.Error(err)
		return
	}
	mp := MatchParser{
		rawHTML: resp,
	}
	matches, err := mp.GetTimeMatches()
	if err != nil {
		t.FailNow()
	}
	if len(matches) > 0 {
		for _, ms := range matches {
			for _, m := range ms {
				logrus.Infof("info: %s", m.GetMDMatchInfo())
			}
		}
	}
}
