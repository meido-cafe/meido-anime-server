package tool

import (
	"errors"
	"strings"
)

type Sql struct {
}

func NewSql() *Sql {
	return &Sql{}
}

type Result struct {
	Sql    string
	Values []any
}

// FormatList 格式化 in 关键词的列表参数
//
// @auth roirea 2022-07-17 22:47:32
// @params
//	n int 参数列表的长度
// @return
//	string	(?,?,?...)格式的sql语句
func (p *Sql) FormatList(n int) string {
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

// CountSql 获取对sql查询结果进行count(*)的sql语句
//
// @auth roirea 2022-08-14 21:56:59
// @params
//	sql string 要进行count(*)的sql语句
// @return
//	string 生成的sql语句
func (p *Sql) CountSql(sql string) string {
	str := `select count(*) from (` + sql + `) count_sql`
	return str
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
func (p *Sql) FormatInsert(tableName string, data map[string]any) (ret Result, err error) {
	if tableName == "" {
		err = errors.New("insert sql 表名为空")
		return
	}
	n := len(data)
	if n == 0 {
		err = errors.New("insert sql 参数为空")
		return
	}

	values := make([]any, 0, n)
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
	sql.WriteString(p.FormatList(n))

	ret.Sql = sql.String()
	ret.Values = values
	return
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
func (p *Sql) FormatInsertBatch(tableName string, list []map[string]any) (ret Result, err error) {
	if tableName == "" {
		err = errors.New("insert sql 表名为空")
		return
	}
	n := len(list) // 参数列表的长度
	if n == 0 {
		err = errors.New("insert list sql 参数为空")
		return
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

	values := make([]any, 0, n*m)
	placeholderAll := make([]string, 0, n) //  [ (?,?,?,?),(?,?,?,?)...]

	for _, obj := range list {
		// 防止for循环map的随机key现象
		for _, field := range fieldList {
			values = append(values, obj[field])
		}
		placeholderAll = append(placeholderAll, placeholder)
	}
	builder.WriteString(strings.Join(placeholderAll, ","))
	ret.Sql = builder.String()
	ret.Values = values
	return
}
