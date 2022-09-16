package model

import (
	"errors"
	"time"
)

type Page struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

func (p *Page) CheckPage() error {
	switch {
	case p.Page <= 0 && p.PageSize <= 0:
		return nil
	case p.Page > 0 && p.PageSize <= 0:
		return errors.New("分页size错误")
	case p.PageSize > 0 && p.Page <= 0:
		return errors.New("分页页码错误")
	}
	return nil
}
func (p *Page) IsPage() bool {
	return p.PageSize > 0 && p.Page > 0
}
func (p *Page) Data() (limit, offset int) {
	return p.PageSize, (p.Page - 1) * p.PageSize
}

type Time struct {
	CreateTime int64 `json:"create_time" db:"create_time"`
	UpdateTime int64 `json:"update_time" db:"update_time"`
}

func NewTime() Time {
	now := time.Now().Unix()
	return Time{
		CreateTime: now,
		UpdateTime: now,
	}
}

type Feed struct {
	Title string     `json:"title"` // rss标题
	Desc  string     `json:"desc"`  // rss描述
	Items []FeedItem `json:"items"` // 种子
}

type FeedItem struct {
	Title   string `json:"title"`    // 种子标题
	Desc    string `json:"desc"`     // 种子描述
	PubDate string `json:"pub_date"` // 发布日期
	Url     string `json:"url"`      // 种子链接
	Length  int    `json:"length"`   // 内容大小, 单位B
}

type UpdatePath struct {
	OldPath string
	NewPath string
}
