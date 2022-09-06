package model

type LinkSource struct {
	Id             int64  // 番剧ID
	Episode        int64  // 集
	SourceFilePath string // 源文件路径
	LinkFilePath   string // 硬链接文件路径
}
