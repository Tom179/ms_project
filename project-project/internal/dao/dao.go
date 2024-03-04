package dao

import (
	"test.com/project-user/internal/database"
	"test.com/project-user/internal/database/gorm"
)

func NewTransaction() *TransactionImpl {
	return &TransactionImpl{
		gorms.NewTran(),
	}
}

type TransactionImpl struct {
	conn database.DbConn
}

func (t TransactionImpl) Action(f func(conn database.DbConn) error) error { //Action的真实执行逻辑
	t.conn.Begin()
	err := f(t.conn) //f是传入的函数名，函数可以在内部手动执行
	if err != nil {
		t.conn.RollBack()
		return err
	}
	t.conn.Commit()
	return nil
}
