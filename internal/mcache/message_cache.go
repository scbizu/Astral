package mcache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/patrickmn/go-cache"
)

const (
	prefix = "message"
)

var messageCache *cache.Cache

func EnableMessageCache() {
	messageCache = cache.New(time.Hour, 2*time.Hour)
}

func IsMessageCacheEnable() bool {
	return messageCache != nil
}

type Hasher interface {
	Hash(string) (string, error)
}

func AddMessage(msg string, h Hasher) error {
	if !IsMessageCacheEnable() {
		return nil
	}
	sum, err := h.Hash(msg)
	if err != nil {
		return err
	}
	if err := messageCache.Add(getMessageKey(sum), "", time.Hour); err != nil {
		return err
	}
	return nil
}

func IsMessageSet(msg string, h Hasher) bool {
	if !IsMessageCacheEnable() {
		return false
	}
	sum, err := h.Hash(msg)
	if err != nil {
		logrus.Errorf("mcache: hash: %q", err)
		return false
	}
	_, ok := messageCache.Get(getMessageKey(sum))
	return ok
}

func getMessageKey(str string) string {
	return fmt.Sprintf("%s-%s", prefix, str)
}

type MD5 struct{}

func (MD5) Hash(msg string) (string, error) {
	hasher := md5.New()
	if _, err := hasher.Write([]byte(msg)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
