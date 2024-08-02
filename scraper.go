package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
	"sync"
)

type IScraper interface {
	Range() []string
	GetProxy() string
	GetCookie() string
	Do(*colly.Collector)
	GetResult() ItemList
}

type BaseScraper struct {
	BaseUrl    string
	SearchList []string
	Proxy      string
	Cookie     string

	Result ItemList
	mutex  sync.Mutex
}

func (s *BaseScraper) Range() (res []string) {
	for _, str := range s.SearchList {
		res = append(res, s.BaseUrl+str)
	}
	return
}

func (s *BaseScraper) GetProxy() string {
	return s.Proxy
}

func (s *BaseScraper) GetCookie() string {
	return s.Cookie
}

func (s *BaseScraper) Do(*colly.Collector) {
	panic("implement me")
}

func (s *BaseScraper) GetResult() ItemList {
	return s.Result
}

type ExHentaiScraper struct {
	BaseScraper
}

func (s *ExHentaiScraper) Do(c *colly.Collector) {
	c2 := c.Clone()
	c2.OnRequest(func(request *colly.Request) {
		request.Headers.Set("cookie", s.GetCookie())
	})
	c2.OnResponse(func(response *colly.Response) {
		url := response.Request.URL.String()
		err := response.Save(basePath + getImageName(url))
		if err != nil {
			log.Printf("save %s error %v", url, err)
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
			_ = c2.Visit(item.Image)
			item.Image = getImageName(item.Image)
			s.mutex.Lock()
			s.Result = append(s.Result, item)
			s.mutex.Unlock()
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
	BaseScraper
}

func (s *EHentaiScraper) Do(c *colly.Collector) {
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
			s.mutex.Lock()
			s.Result = append(s.Result, item)
			s.mutex.Unlock()
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
