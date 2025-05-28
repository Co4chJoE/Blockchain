package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var mysqlLogger logger.Interface

func init() {

	username := "root"  //账号
	password := "root"  //密码
	host := "127.0.0.1" //数据库地址，可以是Ip或者域名
	port := 3306        //数据库端口
	Dbname := "gorm"    //数据库名
	timeout := "10s"    //连接超时，10秒
	// 开启日志功能
	mysqlLogger = logger.Default.LogMode(logger.Info)

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	// 连接成功

	DB = db

}

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	PostCount int    // 用户的文章数量统计字段
	Posts     []Post `gorm:"foreignKey:UserID"` // 一对多关系
}

type Post struct {
	ID           uint `gorm:"primaryKey"`
	Title        string
	Content      string
	UserID       uint      `gorm:"index"`             // 外键，指向 User.ID
	User         User      `gorm:"foreignKey:UserID"` // 关联用户
	CommentCount int       // 文章的评论数量统计字段
	Comments     []Comment `gorm:"foreignKey:PostID"` // 一对多关系
	HasComments  bool      // 如果评论数为0，则设为 false
}

type Comment struct {
	ID      uint `gorm:"primaryKey"`
	Content string
	PostID  uint `gorm:"index"`             // 外键，指向 Post.ID
	Post    Post `gorm:"foreignKey:PostID"` // 关联文章
}

func GetUserPostsWithComments(db *gorm.DB, userID uint) (User, error) {
	var user User
	err := db.Preload("Posts.Comments"). // 二级预加载
						Take(&user, userID).Error
	return user, err
}

func GetMostCommentedPost(db *gorm.DB) (Post, error) {
	var post Post
	err := db.
		Select("posts.*, COUNT(comments.id) AS comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count DESC").
		First(&post).Error
	return post, err
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	// 原子操作更新用户文章数
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + 1")).
		Error
}

func (c *Comment) AfterCreate(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&count).Error; err != nil {
		return err
	}
	updateData := map[string]interface{}{
		"comment_count": count, // 直接使用查询到的准确值
	}
	if count == 0 {
		updateData["has_comments"] = 0
	} else {
		updateData["has_comments"] = 1
	}

	return tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Updates(updateData).Error
}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 查询当前文章的剩余评论数量
	var count int64
	if err := tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&count).Error; err != nil {
		return err
	}

	// 更新文章的评论数量（无论是否为 0）
	updateData := map[string]interface{}{
		"comment_count": count, // 直接使用查询到的准确值
	}

	// 根据评论数量设置状态
	if count == 0 {
		updateData["has_comments"] = 0
	} else {
		updateData["has_comments"] = 1
	}

	// 执行更新操作
	return tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Updates(updateData).Error
}

func main() {
	// DB.AutoMigrate(&User{}, &Post{}, &Comment{})
	// 创建用户(function DONE)
	// user := User{Name: "张三"}
	// DB.Create(&user)
	// fmt.Println(user)
	// fmt.Println(user.ID)
	// 创建文章(function DONE)
	// post := Post{
	// 	Title:  "GORM 使用指南",
	// 	UserID: user.ID, // 正确设置外键关联
	// }
	// result := DB.Create(&post)
	// if result.Error != nil {
	// 	fmt.Println(result.Error)
	// }

	// // 创建评论(function DONE)
	// comments := []*Comment{
	// 	{Content: "好文章！", PostID: post.ID},
	// 	{Content: "非常实用", PostID: post.ID},
	// }
	// result = DB.Create(comments)
	// if result.Error != nil {
	// 	fmt.Println("创建评论失败:", result.Error)
	// 	return
	// }

	// 测试查询
	// user := User{ID: 22}
	// fmt.Println("=== 用户文章及评论 ===")
	// userWithPosts, _ := GetUserPostsWithComments(DB, user.ID)
	// for _, p := range userWithPosts.Posts {
	// 	fmt.Printf("文章《%s》有 %d 条评论\n", p.Title, len(p.Comments))
	// }

	// fmt.Println("=== 最热门文章 ===")
	// hotPost, _ := GetMostCommentedPost(DB)
	// fmt.Printf("最热门文章：%s (评论数：%d)\n", hotPost.Title, hotPost.CommentCount)

	// 测试删除评论
	// fmt.Println("=== 删除评论后的状态 ===")
	// var comment Comment
	// DB.First(&comment)
	// DB.Delete(&comment)

	// var updatedPost Post
	// DB.First(&updatedPost, post.ID)
	// fmt.Printf("文章状态：%d (剩余评论：%d)\n",
	// 	updatedPost.HasComments,
	// 	updatedPost.CommentCount)

	// 查询用户文章及评论(function DONE)
	user, err := GetUserPostsWithComments(DB, 22)
	if err != nil {
		// 处理错误
	}
	fmt.Printf("用户 %s 的文章及评论：%+v\n", user.Name, user.Posts)
	// 查询评论最多的文章(function DONE)
	post, err := GetMostCommentedPost(DB)
	if err != nil {
		// 处理错误
	}
	fmt.Printf("评论最多的文章：%s (评论数：%d)\n", post.Title, post.CommentCount)

	//删除代码测试(function DONE/del hook DONE)
	// var comment Comment
	// DB.Take(&comment, 30)
	// DB.Delete(&comment)
}
