package repo

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
	tool "meido-anime-server/pkg"
	"strings"
)

func NewVideoRepo(db *sqlx.DB) VideoInterface {
	return &VideoRepo{db: db}
}

type VideoRepo struct {
	db *sqlx.DB
}

func (this *VideoRepo) getFields() string {
	fields := []string{
		"id",
		"bangumi_id",
		"title",
		"season",
		"cover",
		"total",
		"rss_url",
		"play_time",
		"create_time",
		"update_time",
	}
	return strings.Join(fields, ",")
}

func (this *VideoRepo) SelectOne(ctx context.Context, id int64, bangumiId int64) (res model.Video, err error) {
	if id <= 0 && bangumiId <= 0 {
		err = errors.New("参数错误")
		log.Println(err)
		return
	}
	q := tool.NewQuery()
	sql := `select ` + this.getFields() + ` from video where true `
	if id > 0 {
		sql += ` and id = ? `
		q.Add(id)
	}
	if bangumiId > 0 {
		sql += ` and bangumi_id = ? `
		q.Add(bangumiId)
	}
	if err = this.db.Get(&res, sql, q.Values()...); err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *VideoRepo) InsertOne(ctx context.Context, video model.Video) (err error) {
	str, values, err := tool.FormatInsert("video", map[string]any{
		"bangumi_id":  video.BangumiId,
		"title":       video.Title,
		"category":    video.Category,
		"origin":      video.Origin,
		"season":      video.Season,
		"cover":       video.Cover,
		"total":       video.Total,
		"rss_url":     video.RssUrl,
		"play_time":   video.PlayTime,
		"create_time": video.CreateTime,
		"update_time": video.UpdateTime,
	})
	if err != nil {
		return err
	}
	if _, err = this.db.Exec(str, values...); err != nil {
		return err
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
	querySQL := `select ` + this.getFields() + ` from video where true`
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

	countSQL := tool.CountSql(querySQL)
	countQuery, err := tx.Query(countSQL, q.Values()...)
	if err != nil {
		return
	}
	defer countQuery.Close()

	if req.Bool() {
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
	sql := ` select ` + this.getFields() + ` from video where id = ? `
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
	sql := ` delete from video where id in ` + tool.FormatList(len(idList))
	_, err = this.db.Exec(sql, q.Values()...)
	return
}

func (this *VideoRepo) UpdateRss(ctx context.Context, id int64, rss string) (err error) {
	sql := ` update video set rss_url = ? where id = ? `
	_, err = this.db.Exec(sql, rss, id)
	return
}
