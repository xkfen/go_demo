package mysql

import (
	"time"
)

type BaseModel struct {
	ID        uint 		`gorm:"primary_key" json:"id" title:"ID"`
	CreatedAt time.Time	`json:"created_at" title:"创建时间"`
	UpdatedAt time.Time	`json:"updated_at" title:"更新时间"`
	DeletedAt *time.Time 	`sql:"index" json:"deleted_at" title:"删除时间"`
}

