package tool

import (
	"github.com/mmcdole/gofeed"
	"meido-anime-server/internal/model"
	"strconv"
)

func ParseRss(url string) (ret model.Feed, err error) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	if err != nil {
		return
	}
	ret.Title = feed.Title
	ret.Desc = feed.Description

	for _, item := range feed.Items {
		length, _ := strconv.Atoi(item.Enclosures[0].Length)
		ret.Items = append(ret.Items, model.FeedItem{
			Title:   item.Title,
			Desc:    item.Description,
			PubDate: item.Published,
			Url:     item.Enclosures[0].URL,
			Length:  length,
		})
	}
	return
}
