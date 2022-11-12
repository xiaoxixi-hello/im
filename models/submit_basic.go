package models

import "gorm.io/gorm"

type Submit struct {
	gorm.Model

	Identity        string `gorm:"column:identity;type:varchar(36);" json:"identity"`                 // 提交表唯一标识
	ProblemIdentity string `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"` // 问题唯一标识
	UserIdentity    string `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`       // 用户唯一标识
	Path            string `gorm:"column:path;type:varchar(255);" json:"path"`                        // 代码路径
}

func (table *Submit) TableName() string {
	return "submit"
}
