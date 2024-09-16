package dao

import (
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/gorm"
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
	err := f(t.conn) //f是传入的函数
	if err != nil {
		t.conn.RollBack()
		return err
	}
	t.conn.Commit()
	return nil
}
