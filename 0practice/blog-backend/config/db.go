package config

import (
	"blog-backend/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	username := "root"  //账号
	password := "root"  //密码
	host := "127.0.0.1" //数据库地址，可以是Ip或者域名
	port := 3306        //数据库端口
	Dbname := "gorm"    //数据库名
	timeout := "10s"    //连接超时，10秒
	// 开启日志功能
	// mysqlLogger = logger.Default.LogMode(logger.Info)

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	// 连接成功
	DB = db
	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

}
