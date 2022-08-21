package repo

import (
	"errors"
	"meido-anime-server/global"
	"strings"
)

// GetExistByMap 获取指定表中指定的字段与值 存在的ID列表
//
// @param
//	tableName string 表名
//	data map[string]interface{} 字段:值 哈希表
// @return
//	ids 结果ID列表
//	error
func GetExistByMap(tableName string, data map[string]any) (ids []int64, err error) {
	if tableName == "" {
		err = errors.New("GetExistByMap 表名为空")
		return
	}
	n := len(data)
	if n == 0 {
		err = errors.New("GetExistByMap 参数为空")
		return
	}
	fields := make([]string, 0, len(data))
	values := make([]any, 0, len(data))
	for k, v := range data {
		fields = append(fields, k)
		values = append(values, v)
	}
	builder := new(strings.Builder)
	builder.WriteString(` select id,`)
	builder.WriteString(strings.Join(fields, ","))
	builder.WriteString(` from `)
	builder.WriteString(tableName)
	builder.WriteString(` where true `)
	for _, item := range fields {
		builder.WriteString(` and `)
		builder.WriteString(item)
		builder.WriteString(` = ? `)
	}
	sql := builder.String()
	ret, err := global.DB.Query(sql, values...)
	if err != nil {
		return
	}
	defer ret.Close()
	var id int64
	for ret.Next() {
		if err = ret.Scan(&id); err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}

// GetNotExistByIdList 检查指定的表中 不存在的id的索引列表
//
//	可用于[批量]更新/删除前检查ID是否存在
func GetNotExistByIdList(tableName string, idList []int64) (indexList []int64, err error) {
	if tableName == "" {
		err = errors.New("GetNotExistByIdList 表名为空")
		return
	}
	n := len(idList)
	if n == 0 {
		err = errors.New("GetNotExistByIdList id为空")
		return
	}
	values := make([]any, 0, n)
	where := make([]string, 0, n)
	for _, item := range idList {
		values = append(values, item)
		where = append(where, "?")
	}
	builder := new(strings.Builder)
	builder.WriteString(` select id,`)
	builder.WriteString(` from `)
	builder.WriteString(tableName)
	builder.WriteString(` where id in ( `)
	builder.WriteString(strings.Join(where, ","))
	builder.WriteString(` ) `)

	sql := builder.String()
	ret, err := global.DB.Query(sql, values...)
	if err != nil {
		return
	}
	defer ret.Close()
	var id int64
	var tmp []int64 // 存在的ID列表
	for ret.Next() {
		if err = ret.Scan(&id); err != nil {
			return
		}
		tmp = append(tmp, id)
	}

	set := make(map[int64]struct{})
	for _, item := range tmp {
		set[item] = struct{}{}
	}
	for index, item := range idList {
		if _, ok := set[item]; !ok {
			indexList = append(indexList, int64(index))
		}
	}
	return
}
