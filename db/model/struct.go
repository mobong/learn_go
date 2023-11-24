package model

type User struct {
	Id    int64
	Name  string  `xorm:"not null index VARCHAR(50) default"`
	Age   int     `xorm:"not null TINYINT(3) default 0"`
	Roles []*Role `xorm:"-"` // 忽略该字段，防止重复查询
}

type Role struct {
	Id   int64
	Name string `xorm:"not null VARCHAR(50) default"`
}
