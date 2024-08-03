package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gocolly/colly/v2"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

type IScraper interface {
	Range() []string
	GetProxy() string
	GetCookie() string
	GetConcurrency() int
	Do(*colly.Collector, string)
	GetUrl(string) string
	GetResult() ItemList
}

type BaseScraper struct {
	BaseUrl    string
	SearchList []string
	Proxy      string
	Cookie     string

	Concurrency int
	MinTime     time.Time

	Result map[string]ItemList
	mutex  sync.Mutex
}

func (s *BaseScraper) Range() (res []string) {
	return s.SearchList
}

func (s *BaseScraper) GetProxy() string {
	return s.Proxy
}

func (s *BaseScraper) GetCookie() string {
	return s.Cookie
}

func (s *BaseScraper) GetConcurrency() int {
	return s.Concurrency
}

func (s *BaseScraper) Do(*colly.Collector, string) {
	panic("implement me")
}

func (s *BaseScraper) GetUrl(key string) string {
	return s.BaseUrl + "?f_search=" + url.QueryEscape(key)
}

func (s *BaseScraper) putItem(key string, item Item) {
	s.mutex.Lock()
	ls := s.Result[key]
	ls = append(ls, item)
	s.Result[key] = ls
	s.mutex.Unlock()
}

func (s *BaseScraper) GetResult() ItemList {
	var ls ItemList
	s.mutex.Lock()
	for _, key := range s.SearchList {
		ll := s.Result[key]
		if len(ll) == 0 {
			log.Printf("search %s not found", key)
		}
		for _, item := range ll {
			if !item.Before(s.MinTime) {
				ls = append(ls, item)
			}
		}
	}
	s.mutex.Unlock()
	return ls
}

type ExHentaiScraper struct {
	*BaseScraper
}

func (s *ExHentaiScraper) Do(c *colly.Collector, key string) {
	c2 := c.Clone()
	c2.OnRequest(func(request *colly.Request) {
		request.Headers.Set("cookie", s.GetCookie())
	})
	c2.OnResponse(func(response *colly.Response) {
		img := response.Request.URL.String()
		err := response.Save(basePath + getImageName(img))
		if err != nil {
			log.Printf("save %s error %v", img, err)
		}
	})

	c.OnHTML("table.itg.gltm", func(table *colly.HTMLElement) {
		table.ForEach("tr", func(i int, tr *colly.HTMLElement) {
			var item, empty Item
			item.Gallery = tr.ChildText("td.gl1m.glcat > div")
			item.Image = tr.ChildAttr("td.gl2m > div.glthumb > div > img", "data-src")
			if item.Image == "" {
				item.Image = tr.ChildAttr("td.gl2m > div.glthumb > div > img", "src")
			}
			item.Date = tr.ChildText("td.gl2m > div:nth-child(3)")
			item.Name = parseName(tr.ChildText("td.gl3m.glname > a > div"))
			item.Link = tr.ChildAttr("td.gl3m.glname > a", "href")
			if item == empty {
				return
			}
			if !item.Before(s.MinTime) {
				_ = c2.Visit(item.Image)
			}
			item.Image = getImageName(item.Image)
			item.Key = key
			s.putItem(key, item)
		})
	})
}

func getImageName(url string) string {
	h := md5.New()
	h.Write([]byte(url))
	bytes := h.Sum(nil)
	return hex.EncodeToString(bytes) + ".jpg"
}

type EHentaiScraper struct {
	*BaseScraper
}

func (s *EHentaiScraper) Do(c *colly.Collector, key string) {
	c.OnHTML("table.itg.gltm", func(table *colly.HTMLElement) {
		table.ForEach("tr", func(i int, tr *colly.HTMLElement) {
			var item, empty Item
			item.Gallery = tr.ChildText("td.gl1m.glcat > div")
			item.Image = tr.ChildAttr("td.gl2m > div.glthumb > div > img", "data-src")
			if item.Image == "" {
				item.Image = tr.ChildAttr("td.gl2m > div.glthumb > div > img", "src")
			}
			item.Date = tr.ChildText("td.gl2m > div:nth-child(3)")
			item.Name = parseName(tr.ChildText("td.gl3m.glname > a > div"))
			item.Link = tr.ChildAttr("td.gl3m.glname > a", "href")
			if item == empty {
				return
			}
			item.Key = key
			s.putItem(key, item)
		})
	})
}

func parseName(s string) string {
	ls := strings.Split(s, "|")
	if len(ls) < 2 {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(ls[len(ls)-1])
}
