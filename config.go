package main

import (
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	ExHentai = Mode("ExHentai")
	EHentai  = Mode("EHentai")
)

type Mode string

func parseMode(s string) Mode {
	reg := regexp.MustCompile("[^a-zA-Z]+")
	ns := reg.ReplaceAllString(s, "")
	if strings.ToLower(ns) == strings.ToLower(string(ExHentai)) {
		return ExHentai
	}
	return EHentai
}

type Config struct {
	Mode       string   `yaml:"Mode"`
	SearchList []string `yaml:"SearchList"`
	Proxy      string   `yaml:"Proxy"`
	Cookie     string   `yaml:"Cookie"`

	TgBotToken string `yaml:"TgBotToken"`
	TgChatId   string `yaml:"TgChatId"`

	Concurrency int `yaml:"Concurrency"`
	MaxDay      int `yaml:"MaxDay"`
}

func (c Config) ToScraper() IScraper {
	bs := &BaseScraper{
		BaseUrl:     "https://e-hentai.org/",
		SearchList:  c.SearchList,
		Proxy:       c.Proxy,
		Cookie:      c.Cookie,
		Concurrency: c.Concurrency,
		MinTime:     time.Now().AddDate(0, 0, -c.MaxDay),
		Result:      make(map[string]ItemList),
		mutex:       sync.Mutex{},
	}
	switch parseMode(c.Mode) {
	case ExHentai:
		bs.BaseUrl = "https://exhentai.org/"
		return &ExHentaiScraper{
			BaseScraper: bs,
		}
	default:
		return &EHentaiScraper{
			BaseScraper: bs,
		}
	}
}
