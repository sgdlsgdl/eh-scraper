package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"golang.org/x/sync/errgroup"
	"log"
	"math"
	"os"
	"sort"
	"time"
)

func main() {
	_ = os.MkdirAll(basePath, os.ModePerm)

	cfg := readConfig(configPath)
	newItems, err := batchFetch(cfg.ToScraper())
	if err != nil {
		log.Printf("batchFetch error %v", err)
		return
	}
	oldItems := readCsv(csvPath)
	total, delta := diff(newItems, oldItems)

	saveCsv(csvPath, total)
	saveMd(mdPath, delta)
	saveHtml(htmlPath, delta)

	batchSend(cfg.Proxy, cfg.TgBotToken, cfg.TgChatId, newItems)
}

type Item struct {
	Gallery string `csv:"Gallery"`
	Image   string `csv:"Image"`
	Date    string `csv:"Date"`
	Name    string `csv:"Name"`
	Link    string `csv:"Link"`

	Ts int64 `csv:"-"`
}

type ItemList []Item

func (ls ItemList) Adjust() {
	for i, v := range ls {
		ts := int64(math.MaxInt64)
		t, err := time.Parse("2006-01-02 15:04", v.Date)
		if err == nil {
			ts = t.Unix()
		}
		ls[i].Ts = ts
	}
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Ts > ls[j].Ts
	})
}

func diff(newItems, oldItems ItemList) (total, delta ItemList) {
	oldM := make(map[string]Item)
	for _, item := range oldItems {
		oldM[item.Link] = item
	}
	for _, item := range newItems {
		if _, ok := oldM[item.Link]; !ok {
			delta = append(delta, item)
		}
	}
	total = append(oldItems, delta...)
	total.Adjust()
	delta.Adjust()
	return
}

func batchFetch(scraper IScraper) (res ItemList, err error) {
	eg := errgroup.Group{}
	for _, url := range scraper.Range() {
		url := url
		eg.Go(func() error {
			return fetch(scraper, url)
		})
	}
	err = eg.Wait()
	res = scraper.GetResult()
	return
}

func fetch(scraper IScraper, url string) error {
	errChan := make(chan error, 1)
	c := colly.NewCollector()
	if scraper.GetProxy() != "" {
		_ = c.SetProxy(scraper.GetProxy())
	}
	if scraper.GetCookie() != "" {
		c.OnRequest(func(request *colly.Request) {
			request.Headers.Set("cookie", scraper.GetCookie())
		})
	}
	scraper.Do(c)
	c.OnError(func(response *colly.Response, err error) {
		putErr(errChan, fmt.Errorf("OnError %s error %v", url, err))
	})
	mErr := c.Visit(url)
	if mErr != nil {
		putErr(errChan, fmt.Errorf("visit %s error %v", url, mErr))
	}
	return getErr(errChan)
}

func putErr(ch chan error, err error) {
	select {
	case ch <- err:
	default:
	}
}

func getErr(ch chan error) error {
	select {
	case err := <-ch:
		return err
	default:
		return nil
	}
}
