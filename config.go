package main

import (
	"regexp"
	"strings"
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
}

func (c Config) ToScraper() IScraper {
	switch parseMode(c.Mode) {
	case ExHentai:
		return &ExHentaiScraper{
			BaseScraper: BaseScraper{
				BaseUrl:    "https://exhentai.org/",
				SearchList: c.SearchList,
				Proxy:      c.Proxy,
				Cookie:     c.Cookie,
			},
		}
	default:
		return &EHentaiScraper{
			BaseScraper: BaseScraper{
				BaseUrl:    "https://e-hentai.org/",
				SearchList: c.SearchList,
				Proxy:      c.Proxy,
				Cookie:     c.Cookie,
			},
		}
	}
}
