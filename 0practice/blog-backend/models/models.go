package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string    `json:"username"  gorm:"unique"`
	Password string    `json:"password" gorm:"column:password"`
	Email    string    `json:"email"`
	Posts    []Post    `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	UserID   uint
	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	UserID  uint
	PostID  uint
	User    User `gorm:"foreignKey:UserID"`
	Post    Post `gorm:"foreignKey:PostID"`
}
