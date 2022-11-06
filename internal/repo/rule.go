package repo

import (
	"github.com/jmoiron/sqlx"
	"log"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/tool"
)

func (this *Repo) RuleSelectList() (res []model.Rule, err error) {
	sql := `
select id,
       name,
       must_contain,
       must_not_contain,
       use_regex,
       episode_filter,
       smart_filter,
       create_time,
       update_time
from rule;
`
	queryx, err := this.db.Queryx(sql)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(queryx *sqlx.Rows) {
		_ = queryx.Close()
	}(queryx)

	for queryx.Next() {
		var v model.Rule
		if err = queryx.StructScan(&v); err != nil {
			log.Println(err)
			return
		}
		res = append(res, v)
	}
	return
}
func (this *Repo) RuleUpdateOne(data model.Rule) (err error) {
	sql := `
update rule
set name=?,
    must_contain=?,
    use_regex=?,
    episode_filter=?,
    smart_filter=?,
    update_time=?
where id = ?;
`
	if _, err = this.db.Exec(sql,
		data.Id,
		data.Name,
		data.MustContain,
		data.MustNotContain,
		data.UseRegex,
		data.EpisodeFilter,
		data.SmartFilter,
		data.UpdateTime,
	); err != nil {
		log.Println(err)
		return
	}

	return
}
func (this *Repo) RuleDeleteOne(id int64) (err error) {
	sql := `
delete from rule where id = ?;
`
	if _, err = this.db.Exec(sql, id); err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *Repo) RuleInsertList(data []model.Rule) (err error) {
	sql := `
insert into rule (name,
                  must_contain,
                  must_not_contain,
                  use_regex,
                  episode_filter,
                  smart_filter,
                  create_time,
                  update_time)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`
	stmt, err := this.db.Prepare(sql)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	for _, item := range data {
		if _, err = stmt.Exec(
			item.Name,
			item.MustContain,
			item.MustNotContain,
			item.UseRegex,
			item.EpisodeFilter,
			item.SmartFilter,
			item.CreateTime,
			item.UpdateTime,
		); err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func (this *Repo) RuleSelectNameListByName(nameList []string) (res []string, err error) {
	q := tool.NewQuery()
	q.Add(nameList)

	sql := `select name from rule where name in` + this.sql.FormatList(len(nameList))
	queryx, err := this.db.Queryx(sql, q.Values()...)
	if err != nil {
		log.Println(err)
		return
	}
	defer queryx.Close()

	for queryx.Next() {
		var name string
		if err = queryx.Scan(&name); err != nil {
			log.Println(err)
			return
		}
		res = append(res, name)
	}
	return
}
