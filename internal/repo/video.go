package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"meido-anime-server/internal/common"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/tool"
	"os"
)

func NewVideoRepo(db *sqlx.DB, sql *tool.Sql, qb *common.QB) VideoInterface {
	return &VideoRepo{db: db, sql: sql, qb: qb}
}

type VideoRepo struct {
	db  *sqlx.DB
	sql *tool.Sql
	qb  *common.QB
}

func (this *VideoRepo) SelectOne(ctx context.Context, id int64, bangumiId int64) (res model.Video, err error) {
	if id <= 0 && bangumiId <= 0 {
		err = errors.New("参数错误")
		log.Println(err)
		return
	}
	q := tool.NewQuery()
	sql := ` 
select 
id,
bangumi_id,
origin,
category,
title,
season,
cover,
total,
rss_url,
play_time,
source_dir,
create_time,
update_time
from video 
where true 
`

	if id > 0 {
		sql += ` and id = ? `
		q.Add(id)
	}
	if bangumiId > 0 {
		sql += ` and bangumi_id = ? `
		q.Add(bangumiId)
	}

	rowx, err := this.db.Queryx(sql, q.Values()...)
	if err != nil {
		log.Println(err)
		return
	}
	if rowx.Next() {
		if err = rowx.StructScan(&res); err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func (this *VideoRepo) InsertOne(ctx context.Context, video model.Video, rule model.QBSetRule) (err error) {
	tx, err := this.db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// DB
	sql, values, err := this.sql.FormatInsert("video", map[string]any{
		"bangumi_id":  video.BangumiId,
		"origin":      video.Origin,
		"category":    video.Category,
		"title":       video.Title,
		"season":      video.Season,
		"cover":       video.Cover,
		"total":       video.Total,
		"rss_url":     video.RssUrl,
		"play_time":   video.PlayTime,
		"source_dir":  video.SourceDir,
		"create_time": video.CreateTime,
		"update_time": video.UpdateTime,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err = tx.Exec(sql, values...); err != nil {
		log.Println(err)
		return err
	}

	/* QB */

	// 添加rss
	ret, err := this.qb.Client.R().SetQueryParams(map[string]string{
		"url":  video.RssUrl,
		"path": rule.RuleName,
	}).Get("/rss/addFeed")
	if err != nil {
		log.Println(err)
		return err
	}
	if ret.IsError() {
		err = errors.New(fmt.Sprintf("[%s] 订阅失败, qbittorrent status: [%d] \n", video.Title, ret.StatusCode))
		log.Println(err)
		return
	}

	ret, err = this.qb.Client.R().SetQueryParams(map[string]string{
		"ruleName": rule.RuleName,
		"ruleDef":  rule.RuleDef,
	}).Get("/rss/setRule")
	if err != nil {
		log.Println(err)
		return
	}

	if ret.IsError() {
		err = errors.New(fmt.Sprintf("[%s] 订阅失败, qbittorrent status: [%d] \n", video.Title, ret.StatusCode))
		log.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *VideoRepo) SelectList(ctx context.Context, req vo.VideoGetListRequest) (res []model.Video, total int64, err error) {
	tx, err := this.db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	q := tool.NewQuery()
	querySQL := `
select 
id,
bangumi_id,
origin,
category,
title,
season,
cover,
total,
rss_url,
play_time,
source_dir,
create_time,
update_time
from video 
where true 
`
	if req.Origin > 0 {
		querySQL += ` and origin = ? `
		q.Add(req.Origin)
	}
	if req.Category > 0 {
		querySQL += ` and category = ? `
		q.Add(req.Category)
	}
	if req.PlayStartTime > 0 {
		querySQL += ` and play_time >= ?`
		q.Add(req.PlayStartTime)
	}
	if req.PlayEndTime > 0 {
		querySQL += ` and play_time <= ?`
		q.Add(req.PlayEndTime)
	}
	if req.AddStartTime > 0 {
		querySQL += ` and create_time >= ?`
		q.Add(req.AddStartTime)
	}
	if req.AddEndTime > 0 {
		querySQL += ` and create_time <= ?`
		q.Add(req.AddEndTime)
	}
	if req.Title != "" {
		querySQL += ` and title like '%?%'`
		q.Add(req.Title)
	}

	countSQL := this.sql.CountSql(querySQL)
	countQuery, err := tx.Query(countSQL, q.Values()...)
	if err != nil {
		return
	}
	defer countQuery.Close()

	if req.IsPage() {
		limit, offset := req.Data()
		q.Add(limit, offset)
		querySQL += ` limit ? offset ? `
	}

	if err = tx.Select(&res, querySQL, q.Values()...); err != nil {
		log.Println(err)
		return
	}

	if countQuery.Next() {
		if err = countQuery.Scan(&total); err != nil {
			log.Println(err)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *VideoRepo) SelectOneById(ctx context.Context, id int64) (ret model.Video, err error) {
	sql := ` 
select 
id,
bangumi_id,
origin,
category,
title,
season,
cover,
total,
rss_url,
play_time,
source_dir,
create_time,
update_time
from video 
where true 
`
	query, err := this.db.Query(sql, id)
	if err != nil {
		return
	}
	defer query.Close()

	if query.Next() {
		err = query.Scan(&ret)
		if err != nil {
			return
		}
	}
	return
}

func (this *VideoRepo) DeleteList(ctx context.Context, idList []int64) (err error) {
	q := tool.NewQuery()
	q.Add(idList)
	sql := ` delete from video where id in ` + this.sql.FormatList(len(idList))
	_, err = this.db.Exec(sql, q.Values()...)
	return
}

func (this *VideoRepo) UpdateRss(ctx context.Context, id int64, rss string) (err error) {
	sql := ` update video set rss_url = ? where id = ? `
	_, err = this.db.Exec(sql, rss, id)
	return
}

func (this *VideoRepo) GetQBLogs(ctx context.Context) (res []model.QBLog, err error) {
	ret, err := this.qb.Client.R().Get("/log/main")
	if err != nil {
		log.Println(err)
		return
	}

	if err = ret.UnmarshalJson(&res); err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *VideoRepo) InsertItemList(ctx context.Context, list []model.VideoItem) (cnt int64, err error) {
	sql := ` insert into video_item (pid,episode,source_path,link_path,create_time,update_time) values (?,?,?,?,?,?)`

	stmt, err := this.db.Preparex(sql)
	if err != nil {
		return
	}

	for _, item := range list {
		// 硬链接
		if err := os.Link(item.SourcePath, item.LinkPath); err != nil {
			log.Printf("[硬链接失败] %s ==> %s [已跳过该集] [error]: %v\n", item.SourcePath, item.LinkPath, err)
			continue
		}

		if _, err := stmt.Exec(item.Pid, item.Episode, item.SourcePath, item.LinkPath, item.CreateTime, item.UpdateTime); err != nil {
			log.Println("插入数据库失败:", err)
			if err = os.Remove(item.LinkPath); err != nil {
				log.Println("删除硬链接文件失败:", err)
			}
		}
		cnt++
	}

	return
}
func (this *VideoRepo) SelectItemList(ctx context.Context, id int64) (res []model.VideoItem, total int64, err error) {
	q := tool.NewQuery()
	sql := `
select 
id,
pid,
episode,
source_path,
link_path,
create_time,
update_time
from video_item
where true 
`
	if id > 0 {
		sql += ` pid = ? `
		q.Add(id)
	}
	queryx, err := this.db.Queryx(sql, q.Values()...)
	if err != nil {
		return
	}
	defer queryx.Close()

	var tmp model.VideoItem
	for queryx.Next() {
		if err = queryx.StructScan(&tmp); err != nil {
			return
		}
		res = append(res, tmp)
	}
	total = int64(len(res))
	return
}
func (this *VideoRepo) DeleteItemList(ctx context.Context, idList []int64) (err error) {
	sql := `
delete from video_item
where true
`
	q := tool.NewQuery()
	if len(idList) > 0 {
		sql += ` and id in ` + this.sql.FormatList(len(idList))
		q.Add(idList)
	}

	if _, err = this.db.Exec(sql, q.Values()...); err != nil {
		return err
	}

	return
}
