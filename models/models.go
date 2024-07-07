package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username  string       `gorm:"type:varchar(50); not null" json:"username"`
	Email     string       `gorm:"type:varchar(100); not null" json:"email"`
	Password  string       `gorm:"type:varchar(255) not null" json:"password"`
	IsRemoved sql.NullBool `gorm:"default:false; not null" json:"isRemoved"`
}

type Articles struct {
	gorm.Model

	Title     string       `gorm:"not null" json:"title"`
	Text      string       `gorm:"type:text not null" json:"text"`
	IsRemoved sql.NullBool `gorm:"default:false; not null" json:"isRemoved"`
	Comments  []Comments   `gorm:"foreignKey:ArticleId"`
	UserId    uint
}

type Comments struct {
	gorm.Model

	Text      string       `gorm:"type:text not null" json:"text"`
	IsRemoved sql.NullBool `gorm:"default:false; not null" json:"isRemoved"`
	UserId    uint         `json:"userId"`
	ArticleId uint         `json:"articleId"`
}
