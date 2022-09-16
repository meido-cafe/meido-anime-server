package repo

import (
	"context"
	"errors"
	"log"
	"meido-anime-server/internal/common"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/tool"
	"os"
)

type VideoSqliteRepo struct {
	db  common.DBClient
	sql *tool.Sql
}

func NewVideoRepo(db common.DBClient, sql *tool.Sql) VideoInterface {
	return &VideoSqliteRepo{
		db:  db,
		sql: sql,
	}
}

func (this *VideoSqliteRepo) CategoryInsertOne(ctx context.Context, cagetory model.Cagetory) (err error) {
	sqlRet, err := this.sql.FormatInsert("category", map[string]any{
		"name":        cagetory.Name,
		"origin":      2,
		"create_time": cagetory.CreateTime,
		"update_time": cagetory.UpdateTime,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err = this.db.Exec(sqlRet.Sql, sqlRet.Values...); err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *VideoSqliteRepo) CategorySelectOne(ctx context.Context, id int64) (ret model.Cagetory, err error) {
	sql := `select id,name,origin,create_time,update_time from category where id = ? `
	rows, err := this.db.Queryx(sql, id)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	if rows.Next() {
		if err = rows.StructScan(&ret); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (this *VideoSqliteRepo) CategorySelectList(ctx context.Context) (ret []model.Cagetory, err error) {
	sql := `select id,name,origin,create_time,update_time from category`
	if err = this.db.Select(&ret, sql); err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *VideoSqliteRepo) CategoryDeleteOne(ctx context.Context, id int64) (err error) {
	sql := `delete from category where id = ? and origin = 2 `
	if _, err = this.db.Exec(sql, id); err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *VideoSqliteRepo) CategoryUpdateName(ctx context.Context, id int64, name string) (err error) {
	sql := ` update category set name = ? where id = ? and origin = 2 `
	if _, err = this.db.Exec(sql, name, id); err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *VideoSqliteRepo) VideoSelectOne(ctx context.Context, id int64, bangumiId int64) (res model.Video, err error) {
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
rule_name,
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
	defer rowx.Close()

	if rowx.Next() {
		if err = rowx.StructScan(&res); err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func (this *VideoSqliteRepo) VideoInsertOne(ctx context.Context, video model.Video) (err error) {
	// DB
	ret, err := this.sql.FormatInsert("video", map[string]any{
		"bangumi_id":  video.BangumiId,
		"origin":      video.Origin,
		"category":    video.Category,
		"title":       video.Title,
		"season":      video.Season,
		"cover":       video.Cover,
		"total":       video.Total,
		"rss_url":     video.RssUrl,
		"rule_name":   video.RuleName,
		"play_time":   video.PlayTime,
		"source_dir":  video.SourceDir,
		"link_dir":    video.LinkDir,
		"create_time": video.CreateTime,
		"update_time": video.UpdateTime,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err = this.db.Exec(ret.Sql, ret.Values...); err != nil {
		log.Println(err)
		return err
	}

	return
}

func (this *VideoSqliteRepo) VideoSelectList(ctx context.Context, request vo.VideoGetListRequest) (res []model.Video, total int64, err error) {
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
rule_name,
play_time,
source_dir,
link_dir,
create_time,
update_time
from video 
where true 
`
	if request.Origin > 0 {
		querySQL += ` and origin = ? `
		q.Add(request.Origin)
	}
	if request.Category > 0 {
		querySQL += ` and category = ? `
		q.Add(request.Category)
	}
	if request.PlayStartTime > 0 {
		querySQL += ` and play_time >= ?`
		q.Add(request.PlayStartTime)
	}
	if request.PlayEndTime > 0 {
		querySQL += ` and play_time <= ?`
		q.Add(request.PlayEndTime)
	}
	if request.AddStartTime > 0 {
		querySQL += ` and create_time >= ?`
		q.Add(request.AddStartTime)
	}
	if request.AddEndTime > 0 {
		querySQL += ` and create_time <= ?`
		q.Add(request.AddEndTime)
	}
	if request.Title != "" {
		querySQL += ` and title like '%?%'`
		q.Add(request.Title)
	}

	countSQL := this.sql.CountSql(querySQL)
	countQuery, err := this.db.Queryx(countSQL, q.Values()...)
	if err != nil {
		log.Println(err)
		return
	}
	defer countQuery.Close()

	if request.IsPage() {
		limit, offset := request.Data()
		q.Add(limit, offset)
		querySQL += ` limit ? offset ? `
	}

	if err = this.db.Select(&res, querySQL, q.Values()...); err != nil {
		log.Println(err)
		return
	}

	if countQuery.Next() {
		if err = countQuery.Scan(&total); err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func (this *VideoSqliteRepo) VideoDeleteList(ctx context.Context, idList []int64) (err error) {
	q := tool.NewQuery()
	q.Add(idList)
	sql := ` delete from video where id in ` + this.sql.FormatList(len(idList))
	if _, err = this.db.Exec(sql, q.Values()...); err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *VideoSqliteRepo) VideoItemInsertList(ctx context.Context, videoItemList []model.VideoItem) (cnt int64, err error) {
	sql := ` insert into video_item (pid,episode,source_path,link_path,create_time,update_time) values (?,?,?,?,?,?)`

	stmt, err := this.db.Prepare(sql)
	if err != nil {
		log.Println(err)
		return
	}

	for _, item := range videoItemList {
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
func (this *VideoSqliteRepo) VideoItemSelectList(ctx context.Context, pid []int64) (res []model.VideoItem, total int64, err error) {
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
	n := len(pid)
	if n > 0 {
		sql += ` and pid in ` + this.sql.FormatList(n)
		q.Add(pid)
	}

	if err = this.db.Select(&res, sql, q.Values()...); err != nil {
		log.Println(err)
		return
	}

	total = int64(len(res))
	return
}

func (this *VideoSqliteRepo) VideoItemDeleteList(ctx context.Context, idList []int64) (err error) {
	sql := `
delete from video_item
where id in ` + this.sql.FormatList(len(idList))

	q := tool.NewQuery()
	q.Add(idList)

	if _, err = this.db.Exec(sql, q.Values()...); err != nil {
		log.Println(err)
		return err
	}

	return
}

func (this *VideoSqliteRepo) VideoItemUpdateLinkPathList(ctx context.Context, videoItemList []model.VideoItem) (err error) {
	sql := ` update video_item set link_path = ? where id = ? `
	stmt, err := this.db.Prepare(sql)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	for _, item := range videoItemList {
		if _, err = stmt.Exec(item.LinkPath, item.Id); err != nil {
			log.Println(err)
			return err
		}
	}
	return
}

func (this *VideoSqliteRepo) VideoItemDeleteListByPid(ctx context.Context, pid []int64) (err error) {
	q := tool.NewQuery()
	sql := `delete from video_item where pid in ` + this.sql.FormatList(len(pid))
	q.Add(pid)
	if _, err := this.db.Exec(sql, q.Values()...); err != nil {
		log.Println(err)
		return err
	}
	return
}

func (this *VideoSqliteRepo) VideoUpdateCategoryList(ctx context.Context, videoList []model.Video) (err error) {
	sql := ` update video set category = ?, link_dir = ? where id = ? `
	stmt, err := this.db.Prepare(sql)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	for _, item := range videoList {
		if _, err = stmt.Exec(item.Category, item.LinkDir, item.Id); err != nil {
			log.Println(err)
			return err
		}
	}
	return
}

func (this *VideoSqliteRepo) CategorySelectOneByName(ctx context.Context, name string) (ret model.Cagetory, err error) {
	sql := ` select id,name,origin,create_time,update_time from category where name = ? `
	rows, err := this.db.Queryx(sql, name)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	if rows.Next() {
		if err = rows.StructScan(&ret); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (this *VideoSqliteRepo) VideoSelectListById(ctx context.Context, idList []int64) (ret []model.Video, err error) {
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
rule_name,
play_time,
source_dir,
link_dir,
create_time,
update_time
from video 
where id in 
` + this.sql.FormatList(len(idList))
	q := tool.NewQuery()
	q.Add(idList)

	queryx, err := this.db.Queryx(sql, q.Values()...)
	if err != nil {
		log.Println(err)
		return
	}
	defer queryx.Close()

	var v model.Video
	for queryx.Next() {
		if err = queryx.StructScan(&v); err != nil {
			log.Println(err)
			return
		}
		ret = append(ret, v)
	}
	return
}
