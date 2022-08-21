package pkg

import (
	"errors"
	"meido-anime-server/internal/enum"
	"regexp"
	"strconv"
	"strings"
)

// PreHandleTitle 预处理
func PreHandleTitle(title string) string {
	// 替换所有中文中括号为英文中括号
	title = strings.ReplaceAll(title, "【", "[")
	title = strings.ReplaceAll(title, "】", "]")

	// 去除中括号内周围的空格
	pre, _ := regexp.Compile("\\[\\s+")
	title = pre.ReplaceAllString(title, "[")
	suffix, _ := regexp.Compile("\\s+\\]")
	title = suffix.ReplaceAllString(title, "]")
	return title
}

// GetSeason 提取季
//
// S02,S2,s2,s02,第2季,第2期,第二季,第二期
func GetSeason(title string) (season int64, matchStr string, err error) {
	var ret string
	var reg *regexp.Regexp

	// 中文数字
	reg, err = regexp.Compile("\\s第[\u4e00\u4e8c\u4e09\u56db\u4e94\u516d\u4e03\u516b\u4e5d\u5341]+季|\\s第[\u4e00\u4e8c\u4e09\u56db\u4e94\u516d\u4e03\u516b\u4e5d\u5341]+期")
	if err != nil {
		return
	}
	ret = reg.FindString(title)
	if ret != "" {
		matchStr = ret
		ret = strings.TrimSpace(ret)
		ret = strings.Replace(ret, "第", "", 1)
		ret = strings.Replace(ret, "季", "", 1)
		ret = strings.Replace(ret, "期", "", 1)
		season, err = handleZhNumber(ret)
		return
	}

	// 数学数字
	reg, err = regexp.Compile("s\\d+|S\\d+|第?\\d+[季期]")
	if err != nil {
		return
	}
	ret = reg.FindString(title)
	if ret != "" {
		matchStr = ret
		reg, err = regexp.Compile("\\d+")
		if err != nil {
			return
		}
		ret = reg.FindString(ret)
		var seasonTmp int
		seasonTmp, err = strconv.Atoi(ret)
		season = int64(seasonTmp)
		return
	}

	// 罗马数字
	reg, err = regexp.Compile("\\s第?[IVXLCDM]+")
	if err != nil {
		return
	}
	ret = reg.FindString(title)
	if ret != "" {
		matchStr = ret
		ret = strings.TrimSpace(ret)
		ret = strings.Replace(ret, "第", "", 1)
		season = Rome2int(ret)
		return
	}

	return
}

// GetEpisode 提取集
//
// [5],[05], 5 , 05 , 第三十六集, 第05集 , 第三十六话, 第05话
func GetEpisode(title string) (episode int64, err error) {
	var reg *regexp.Regexp
	var tmp int

	// 数字集数
	reg, err = regexp.Compile("\\[\\d+\\]|\\s\\d+\\s|第\\d+话|第\\d+集")
	if err != nil {
		return
	}
	ret := reg.FindString(title)
	if ret != "" {
		reg, err = regexp.Compile("\\d+")
		if err != nil {
			return
		}

		tmp, err = strconv.Atoi(reg.FindString(ret))
		if err != nil {
			return
		}
		episode = int64(tmp)
		return
	}

	// 尝试提取中文集数
	reg, err = regexp.Compile("第[\u4e00\u4e8c\u4e09\u56db\u4e94\u516d\u4e03\u516b\u4e5d\u5341]+话|第[\u4e00\u4e8c\u4e09\u56db\u4e94\u516d\u4e03\u516b\u4e5d\u5341]+集")
	if err != nil {
		return
	}
	ret = reg.FindString(title)
	if ret != "" {
		ret = strings.TrimSpace(ret)
		ret = strings.Replace(ret, "第", "", 1)
		ret = strings.Replace(ret, "集", "", 1)
		ret = strings.Replace(ret, "话", "", 1)
		episode, err = handleZhNumber(ret)
		return
	}
	err = errors.New("不支持的集数格式")
	return
}

func Rome2int(rome string) (ret int64) {
	tmp := 0
	m := map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	last := 0
	for i := len(rome) - 1; i >= 0; i-- {
		temp := m[rome[i]]
		sign := 1
		if temp < last {
			sign = -1
		}

		tmp += sign * temp
		last = temp
	}
	ret = int64(tmp)
	return
}

func handleZhNumber(str string) (season int64, err error) {
	arr := []rune(str)
	switch len(arr) {
	case 0:
		err = errors.New("格式错误")
		return
	case 1:
		if arr[0] == '十' {
			return 10, nil
		} else {
			return enum.NumberMap[arr[0]], nil
		}
	case 2:
		if arr[0] == '十' {
			return 10 + enum.NumberMap[arr[1]], nil
		}
		if arr[1] == '十' {
			return enum.NumberMap[arr[0]] * 10, nil
		}
		err = errors.New("格式错误")
		return
	case 3:
		if arr[0] == '十' || arr[1] != '十' || arr[2] == '十' {
			err = errors.New("番剧季信息格式错误")
			return
		}
		return enum.NumberMap[arr[0]]*10 + enum.NumberMap[arr[2]], nil
	default:
		err = errors.New("不支持的格式")
		return
	}
}
