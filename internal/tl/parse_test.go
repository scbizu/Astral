package tl

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestMatchParser_GetTimelines(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	r, err := ioutil.ReadFile("test.json")
	if err != nil {
		t.Fatal(err)
	}
	_, resp, err := newParseRespFromReader(bytes.NewBuffer(r))
	if err != nil {
		t.Fatal(err)
	}
	mp := MatchParser{
		rawHTML: resp,
	}
	tls, err := mp.GetTimelines()
	if err != nil {
		t.Fatal(err)
		return
	}
	if len(tls) == 0 {
		t.Fatalf("parser: expected matches, found no matches")
	}
	logrus.Infof("timelines: %#v", tls)
}

func TestMatchParser_GetTimeMatches(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	r, err := ioutil.ReadFile("test.json")
	if err != nil {
		t.Error(err)
		return
	}
	_, resp, err := newParseRespFromReader(bytes.NewBuffer(r))
	if err != nil {
		t.Error(err)
		return
	}
	mp := MatchParser{
		rawHTML: resp,
	}
	matches, err := mp.GetTimeMatches()
	if err != nil {
		t.Error(err)
		return
	}
	if len(matches) == 0 {
		t.Error("parser: expected matches, found no matches")
	}
	if len(matches) > 0 {
		for _, ms := range matches {
			for _, m := range ms {
				logrus.Infof("info: %s", m.GetMDMatchInfo())
			}
		}
	}
}

func TestParseTS(t *testing.T) {
	ti, err := time.Parse(timeFmt, "July 11, 2019 - 11:00 UTC")
	if err != nil {
		t.Fatal(err)
	}
	cn, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}
	countDown := time.Until(ti.In(cn))
	t.Log(int(countDown))
}
