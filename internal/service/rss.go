package service

import (
	"github.com/gocolly/colly/v2"
	"log"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/repo"
	"meido-anime-server/pkg"
	"net/url"
	"strconv"
	"strings"
)

type RssService struct {
	repo repo.RssInterface
}

func NewRssService(repo repo.RssInterface) *RssService {
	return &RssService{repo: repo}
}

func (this *RssService) GetInfoMikan(request vo.GetRssInfoMikanRequest) (response vo.GetRssMiaknInfoResponse, err error) {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"), colly.MaxDepth(1))

	// 获取番剧ID
	c.OnHTML("div[id='sk-container'] > div[class='central-container'] > ul ", func(e *colly.HTMLElement) {
		href := e.ChildAttr("li > a", "href")
		{
			arr := strings.Split(href, "/")
			midString := arr[len(arr)-1]
			mid, _ := strconv.Atoi(midString)
			response.Mid = int64(mid)
		}
	})

	// 获取所有字幕组信息
	c.OnHTML("div[class='leftbar-nav'] > ul[class='list-unstyled']", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, element *colly.HTMLElement) {
			name := element.ChildText("span > a")
			gidString := element.ChildAttr("span > a", "data-subgroupid")
			if gidString != "" {
				gid, _ := strconv.Atoi(gidString)
				response.Group = append(response.Group, vo.MikanGroup{
					Gid:       int64(gid),
					GroupName: name,
				})
			}
		})
	})

	if err = c.Visit("https://mikanani.me/Home/Search?searchstr=" + url.QueryEscape(request.SubjectName)); err != nil {
		log.Println(err)
		return
	}
	return
}

func (this *RssService) GetSearch(request vo.GetRssSearchRequest) (response vo.GetRssSearchResponse, err error) {
	response.Url = "https://mikanani.me/RSS/Search?searchstr=" + url.QueryEscape(request.SubjectName)
	response.Feed, err = pkg.ParseRss(response.Url)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (this *RssService) GetSubject(request vo.GetRssSubjectRequest) (response vo.GetRssSubjectResponse, err error) {
	response.Url = "https://mikanani.me/RSS/Bangumi?bangumiId=" + strconv.Itoa(int(request.MikanId)) + "&subgroupid=" + strconv.Itoa(int(request.MikanGroupId))
	response.Feed, err = pkg.ParseRss(response.Url)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
