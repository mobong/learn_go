package main

import (
	"fmt"
	"learn_go/db/model"
	"xorm.io/xorm"
)

var db *xorm.Engine

func main() {
	syncSchema()
	userOperate()
	defer db.Close()
}

// 同步创建user表
func syncSchema() {
	db.StoreEngine("InnoDB").Sync2(
		new(model.User),
		new(model.Role),
	)
}

func userOperate() {
	// 插入数据
	user := model.User{Name: "John Doe", Age: 30}
	affected, err := db.Insert(&user)
	if err != nil {
		// 处理插入错误
		fmt.Println(err)
	} else {
		fmt.Println("Rows Affected:", affected)
		fmt.Println("Last Insert ID:", user.Id)
	}

	// 更新数据
	user = model.User{Name: "Jane Doe", Age: 31}
	affected, err = db.ID(1).Update(&user)
	if err != nil {
		// 处理更新错误
		fmt.Println(err)
	} else {
		fmt.Println("Rows Affected:", affected)
	}

	// 查询单行数据
	user = model.User{}
	has, err := db.ID(1).Get(&user)
	if err != nil {
		// 处理查询错误
		fmt.Println(err)
	} else if has {
		fmt.Println("User:", user)
	} else {
		fmt.Println("User not found")
	}

	// 查询多行数据: 查询前2条年龄大于等于18岁的用户，按照年龄降序排序
	users := make([]model.User, 0)
	err = db.Where("age >= ?", 18).OrderBy("age DESC").Limit(2).Find(&users)
	if err != nil {
		// 处理查询错误
		fmt.Println(err)
	} else {
		fmt.Println("Users:", users)
	}

	// 多表关联查询用户及其关联的所有角色
	role := model.Role{Name: "John Doe"}
	db.Insert(&role)
	user = model.User{}
	has, err = db.ID(1).Get(&user)
	if err != nil {
		// 处理查询错误
		fmt.Println(err)
	} else if has {
		err := db.Table("user").Alias("u").
			Join("INNER", "role", "u.id = role.id").
			Where("u.id = ?", user.Id).
			Find(&user.Roles)
		if err != nil {
			// 处理查询错误
			fmt.Println(err)
		} else {
			fmt.Println("User:", user)
			fmt.Println("Roles:", user.Roles)
		}
	} else {
		fmt.Println("User not found")
	}

	// 删除数据
	user = model.User{Id: 1}
	affected, err = db.Delete(&user)
	if err != nil {
		// 处理删除错误
		fmt.Println(err)
	} else {
		fmt.Println("Rows Affected:", affected)
	}
}

func init() {
	database, err := dbStartUp()
	if err != nil {
		panic(err)
	}
	db = database
}
