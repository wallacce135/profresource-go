package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username string     `gorm:"type:varchar(50) not null" json:"username"`
	Email    string     `gorm:"type:varchar(100) not null" json:"email"`
	Password string     `gorm:"type:varchar(255) not null" json:"password"`
	Articles []Articles `gorm:"foreignKey:UserId constaint:onUpdate:CASCADE,onDelete:SET NULL"`
	Comments []Comments `gorm:"foreignKey:UserId constraint:onUpdate:CASCADE,onDelete:CASCADE"`
}

type Articles struct {
	gorm.Model

	Title    string     `gorm:"not null" json:"title"`
	Text     string     `gorm:"type:text not null" json:"text"`
	Comments []Comments `gorm:"foreignKey:ArticleId constraint:OnUpdate:CASCADE,onDelete:CASCADE"`
	UserId   uint
}

type Comments struct {
	gorm.Model

	Text      string `gorm:"type:text not null" json:"text"`
	UserId    uint
	ArticleId uint
}
