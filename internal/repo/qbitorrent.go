package repo

import (
	"context"
	"fmt"
	"log"
	"meido-anime-server/internal/common"
	"meido-anime-server/internal/global"
	"strings"
)

type Qbittorrent struct {
	qb *common.QB
}

func NewQbittorrent() *Qbittorrent {
	return &Qbittorrent{qb: common.GetQB()}
}

func (q Qbittorrent) DeleteTorrent(ctx context.Context, hashes []string, deleteFile bool) (err error) {
	ret, err := q.qb.Client.R().SetQueryParamsAnyType(map[string]any{
		"hashes":      strings.Join(hashes, "|"),
		"deleteFiles": deleteFile,
	}).Get("/rss/setRule")

	if err != nil {
		log.Println(err)
		return err
	}
	if ret.IsError() {
		err = fmt.Errorf("%s", ret.String())
		log.Println(err)
		return
	}
	return
}

func (q Qbittorrent) SetRule(ctx context.Context, data map[string]string) (err error) {
	ret, err := q.qb.Client.R().SetQueryParams(data).Get("/rss/setRule")
	if err != nil {
		log.Println(err)
		return err
	}
	if ret.IsError() {
		err = fmt.Errorf("%s", ret.String())
		log.Println(err)
		return
	}
	return
}
func (q Qbittorrent) DeleteRule(ctx context.Context, ruleName string) (err error) {
	ret, err := q.qb.Client.R().SetQueryParams(map[string]string{
		"ruleName": ruleName,
	}).Get("/rss/removeRule")
	if err != nil {
		log.Println(err)
		return err
	}
	if ret.IsError() {
		err = fmt.Errorf("%s", ret.String())
		log.Println(err)
		return
	}
	return
}
func (q Qbittorrent) AddRss(ctx context.Context, rss string, name string) (err error) {
	ret, err := q.qb.Client.R().SetQueryParams(map[string]string{
		"url":  rss,
		"path": global.RssFolder + "/" + name,
	}).Get("/rss/addFeed")
	if err != nil {
		log.Println("qbittorrent订阅rss失败: ", err)
		return err
	}
	if ret.IsError() {
		err = fmt.Errorf("%s", ret.String())
		log.Println("qbittorrent订阅rss失败: ", err)
		return
	}
	return
}

func (q Qbittorrent) DeleteRss(ctx context.Context, rss string) (err error) {
	ret, err := q.qb.Client.R().SetQueryParams(map[string]string{
		"path": rss,
	}).Get("/rss/removeItem")
	if err != nil {
		log.Println(err)
		return err
	}
	if ret.IsError() {
		err = fmt.Errorf("%s", ret.String())
		log.Println(err)
		return
	}
	return
}

func (q Qbittorrent) GetLog(ctx context.Context) (res []map[string]any, err error) {
	ret, err := q.qb.Client.R().Get("/log/main")
	if err != nil {
		log.Println(err)
		return
	}
	if ret.IsError() {
		err = fmt.Errorf("%s", ret.String())
		log.Println(err)
		return
	}

	if err = ret.Unmarshal(&res); err != nil {
		log.Println(err)
		return
	}
	return
}

func (q Qbittorrent) AddTorrents(ctx context.Context, torrentList []string, categoryName string, savePath string) (err error) {
	ret, err := q.qb.Client.R().SetFormDataAnyType(map[string]interface{}{
		"urls":     torrentList,
		"category": categoryName,
		"savepath": savePath,
	}).Post("/torrents/add")
	if err != nil {
		log.Println(err)
		return err
	}
	if ret.IsError() {
		err = fmt.Errorf("%s", ret.String())
		log.Println(err)
		return
	}
	return
}
