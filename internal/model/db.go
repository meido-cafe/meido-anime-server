package model

type Video struct {
	Id        int64  `db:"id"`
	BangumiId int64  `db:"bangumi_id"`
	Title     string `db:"title"`
	Cover     string `db:"cover"`
	Total     int64  `db:"total"`
	RssUrl    string `db:"rss_url"`
	Time
}

type VideoItem struct {
	Id          int64  `db:"id"`
	TorrentHash string `db:"torrent_hash"`
	Pid         int64  `db:"pid"`
	Status      int64  `db:"status"`
	Season      int64  `db:"season"`
	Episode     int64  `db:"episode"`
	SavePath    string `db:"save_path"`
	Time
}
