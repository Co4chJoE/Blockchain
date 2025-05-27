package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL 驱动（匿名导入）
	"github.com/jmoiron/sqlx"
)

// 定义与数据库表对应的结构体
type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

func getHighestPaidEmployee(db *sqlx.DB) (Employee, error) {
	var emp Employee
	err := db.Get(&emp, `
        SELECT id, name, department, salary 
        FROM employees 
        ORDER BY salary DESC 
        LIMIT 1`)

	if err != nil {
		return Employee{}, err
	}
	return emp, nil
}
func main() {
	// 假设已经建立数据库连接
	db, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/gorm?charset=utf8")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 查询技术部员工
	var employees []Employee
	err = db.Select(&employees, `
        SELECT id, name, department, salary 
        FROM employees 
        WHERE department = ?`, "技术部")

	if err != nil {
		panic(err)
	}

	fmt.Printf("技术部员工共 %d 人：\n", len(employees))
	for _, emp := range employees {
		fmt.Printf("%+v\n", emp)
	}
	highestPaid, err := getHighestPaidEmployee(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("工资最高的员工：")
	fmt.Printf("%+v\n", highestPaid)

}
