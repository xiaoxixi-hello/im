package models

import "gorm.io/gorm"

type Problem struct {
	gorm.Model        // 默认的id 创建、更新、删除时间等
	Identity   string `gorm:"column:identity;type:varchar(36);" json:"identity"`        // 问题表唯一标识
	CategoryId string `gorm:"column:category_id;type:varchar(255);" json:"category_id"` // 分类id, 以逗号分割
	Title      string `gorm:"column:title;type:varchar(255);" json:"title"`             // 文章标题
	Content    string `gorm:"column:content;type:text;" json:"content"`                 // 文章正文
	MaxRuntime int    `gorm:"column:max_runtime;type:int(11);" json:"max_runtime"`      // 最大运行时长
	MaxMem     int    `gorm:"column:max_mem;type:int(11);" json:"max_mem"'`             // 最大运行内存
}

func (table *Problem) TableName() string {
	return "problem"
}
