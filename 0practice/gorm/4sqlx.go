package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

func main() {
	// 1. 连接数据库（请替换实际数据库凭据）
	db, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/gorm?charset=utf8")
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 2. 执行复杂查询：价格 > 50 元
	var books []Book
	const minPrice = 50.0
	err = db.Select(&books, `
        SELECT id, title, author, price 
        FROM books 
        WHERE price > ?
        ORDER BY price DESC`, minPrice)

	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	// 3. 输出结果
	fmt.Printf("找到 %d 本价格超过 %.2f 元的书籍:\n", len(books), minPrice)
	for _, book := range books {
		fmt.Printf("ID: %d | 书名: %-20s | 作者: %-10s | 价格: ￥%.2f\n",
			book.ID, book.Title, book.Author, book.Price)
	}
}
