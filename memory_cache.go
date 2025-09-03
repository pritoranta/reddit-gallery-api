package main

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

// 5 day cache, cleaned once per hour
var c = cache.New(120*time.Hour, 1*time.Hour)

func setSubredditsToCache(q subredditQuery, s []subreddit) {
	c.Set(fmt.Sprintf("%+v", q), s, cache.DefaultExpiration)
}

func getCachedSubreddits(q subredditQuery) (s []subreddit, found bool) {
	subreddits, found := c.Get(fmt.Sprintf("%+v", q))
	if found {
		return subreddits.([]subreddit), true
	}
	return nil, false
}

func setMediaToCache(q mediaQuery, m mediaList) {
	c.Set(fmt.Sprintf("%+v", q), m, cache.DefaultExpiration)
}

func getCachedMedia(q mediaQuery) (m mediaList, found bool) {
	media, found := c.Get(fmt.Sprintf("%+v", q))
	if found {
		return media.(mediaList), true
	}
	return mediaList{}, false
}
