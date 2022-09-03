package repo

import (
	"github.com/jmoiron/sqlx"
	"log"
	"meido-anime-server/internal/model"
)

type CronRepo struct {
	db *sqlx.DB
}

func NewCronRepo(db *sqlx.DB) *CronRepo {
	return &CronRepo{db: db}
}

func (this *CronRepo) GetVideoList() (res []model.Video, err error) {
	sql := `select  id, bangumi_id, title, season, cover, total, rss_url, play_time, download_path,create_time, update_time from video `
	queryx, err := this.db.Queryx(sql)
	if err != nil {
		log.Println(err)
		return
	}
	defer queryx.Close()

	var tmp model.Video
	for queryx.Next() {
		if err = queryx.StructScan(&tmp); err != nil {
			log.Println(err)
		}
		res = append(res, tmp)
	}
	return
}
