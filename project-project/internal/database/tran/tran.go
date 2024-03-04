package tran

import "test.com/project-user/internal/database"

type Transaction interface { //事务操作需要数据库连接
	Action(func(conn database.DbConn) error) error //传入参数是一个函数
}
