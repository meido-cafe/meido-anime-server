package pkg

import (
	"errors"
	"strings"
)

// FormatList 格式化 in 关键词的列表参数
//
// @auth roirea 2022-07-17 22:47:32
// @params
//	n int 参数列表的长度
// @return
//	string	(?,?,?...)格式的sql语句
func FormatList(n int) string {
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		arr[i] = "?"
	}
	s := strings.Join(arr, ",")

	sql := new(strings.Builder)
	sql.WriteString(" ( ")
	sql.WriteString(s)
	sql.WriteString(" ) ")
	return sql.String()
}

// FormatConflict 格式化 conflict sql语句
//
// @auth roirea 2022-07-17 22:40:58
// @param
//	key []string 唯一约束键
//	field []string 重复则更新的字段
// @return
//	string 生成的conflict sql语句
//	error
func FormatConflict(key []string, field []string) (string, error) {
	if len(key) <= 0 {
		return "", errors.New("唯一约束键为空")
	}
	if len(field) <= 0 {
		return "", errors.New("更新字段为空")
	}
	builder := new(strings.Builder)
	builder.WriteString(" on conflict ( ")
	builder.WriteString(strings.Join(key, ","))
	builder.WriteString(" ) do update set ")

	for i, item := range field {
		builder.WriteString(item)
		builder.WriteString("=excluded.")
		builder.WriteString(item)
		if i < len(field)-1 {
			builder.WriteString(",")
		}
	}
	return builder.String(), nil
}

// FormatInsert 格式化插入sql语句
//
// @auth roirea 2022-07-17 22:33:57
// @param
//	tableName string 表名
//	data	map[string]interface{} key为字段名,value为值的哈希表
// @return
//	string 生成的sql语句
//	[]interface{} 执行参数
//	error
func FormatInsert(tableName string, data map[string]interface{}) (string, []interface{}, error) {
	if tableName == "" {
		return "", nil, errors.New("insert sql 表名为空")
	}
	n := len(data)
	if n == 0 {
		return "", nil, errors.New("insert sql 参数为空")
	}

	values := make([]interface{}, 0, n)
	fields := make([]string, 0, n)

	sql := new(strings.Builder)
	sql.WriteString(" insert into ")
	sql.WriteString(tableName)
	sql.WriteString(" ( ")

	for k, v := range data {
		fields = append(fields, k)
		values = append(values, v)
	}

	sql.WriteString(strings.Join(fields, ","))
	sql.WriteString(" ) values ")
	sql.WriteString(FormatList(n))
	return sql.String(), values, nil
}

// FormatInsertBatch 构造批量插入SQL语句
//
// @auth roirea 2022-07-17 22:42:11
// @params
//	tableName string	表名
//	list []map[string]interface{}	参数列表,key为字段名,value为值
// @return
//	string 生成的sql语句
//	[]interface{} 执行参数
//	error
func FormatInsertBatch(tableName string, list []map[string]interface{}) (string, []interface{}, error) {
	if tableName == "" {
		return "", nil, errors.New("insert sql 表名为空")
	}
	n := len(list) // 参数列表的长度
	if n == 0 {
		return "", nil, errors.New("insert list sql 参数为空")
	}
	m := len(list[0]) // 字段数量

	builder := new(strings.Builder)

	// 获取字段列表
	fieldList := make([]string, 0, m)
	for k, _ := range list[0] {
		fieldList = append(fieldList, k)
	}
	builder.WriteString(" (")
	builder.WriteString(strings.Join(fieldList, ","))
	builder.WriteString(") ")
	fields := builder.String()
	builder.Reset()

	// 构建参数占位符字符串
	placeholderList := make([]string, m)
	for i := 0; i < m; i++ {
		placeholderList[i] = "?"
	}
	builder.WriteString(" (")
	builder.WriteString(strings.Join(placeholderList, ","))
	builder.WriteString(") ")
	placeholder := builder.String()
	builder.Reset()

	// fieldList : (f1,f2,f3,f4,f5)
	// placeholder : (?,?,?,?,?)
	builder.WriteString(" insert into ")
	builder.WriteString(tableName)
	builder.WriteString(fields)
	builder.WriteString(" values ")

	values := make([]interface{}, 0, n*m)
	placeholderAll := make([]string, 0, n) //  [ (?,?,?,?),(?,?,?,?)...]

	for _, obj := range list {
		// 防止for循环map的随机key现象
		for _, field := range fieldList {
			values = append(values, obj[field])
		}
		placeholderAll = append(placeholderAll, placeholder)
	}
	builder.WriteString(strings.Join(placeholderAll, ","))
	return builder.String(), values, nil
}

// FormatInsertConflict 构造存在就更新的插入SQL语句
//
// @auth roirea 2022-07-17 22:44:47
// @params
//	tableName string	表名
//	data map[string]interface{}	key为字段名,value为值的参数哈希表
//	key []string	on conflict 约束键
//	updateList []string	要更新的字段列表, 字段1=excluded.字段1
// @return
//	string	生成的sql语句
//	[]interface{} 用于执行的参数列表
//	error
func FormatInsertConflict(tableName string, data map[string]interface{}, key []string, updateList []string) (string, []interface{}, error) {
	sql, values, err := FormatInsert(tableName, data)
	if err != nil {
		return "", nil, err
	}

	csql, err := FormatConflict(key, updateList)
	if err != nil {
		return "", nil, err
	}

	return sql + csql, values, nil
}

// FormatInsertBatchConflict 构造批量存在就更新的插入SQL语句
//
// @auth roirea 2022-07-17 22:46:00
// @params
//	tableName string	表名
//	list []map[string]interface{}	参数列表,key为字段名,value为值
//	key []string	on conflict 约束键
//	updateList []string	要更新的字段列表, 字段1=excluded.字段1
// @return
//	string	生成的sql语句
//	[]interface{} 用于执行的参数列表
//	error
func FormatInsertBatchConflict(tableName string, list []map[string]interface{}, key []string, updateList []string) (string, []interface{}, error) {
	sql, value, err := FormatInsertBatch(tableName, list)
	if err != nil {
		return "", nil, err
	}
	csql, err := FormatConflict(key, updateList)
	if err != nil {
		return "", nil, err
	}
	return sql + csql, value, err
}
