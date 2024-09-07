package gorms //获取数据库连接

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"test.com/project-project/config"
)

var _db *gorm.DB

func init() {
	username := config.C.MysqlConfig.Username
	password := config.C.MysqlConfig.Password
	host := config.C.MysqlConfig.Host
	port := config.C.MysqlConfig.Port
	Dbname := config.C.MysqlConfig.Db
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	var err error
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	//_db.AutoMigrate(&member.Member{}, &organization.Organization{}) //自动建表
}

func GetDB() *gorm.DB {
	return _db
}

type GormConn struct { //存放数据库连接的结构体
	db *gorm.DB
	tx *gorm.DB
}

func New() *GormConn { //获取数据库连接
	return &GormConn{db: GetDB()}
}

func NewTran() *GormConn {
	return &GormConn{db: GetDB(), tx: GetDB()}
}

func (g *GormConn) Session(ctx context.Context) *gorm.DB { //新建session,为什么要单独封装结构体来写session方法？
	return g.db.Session(&gorm.Session{Context: ctx})
}

func (g *GormConn) Tx(ctx context.Context) *gorm.DB {
	return g.tx.WithContext(ctx)
}

func (g *GormConn) Begin() { //每次事务都要重新获取新的连接
	g.tx = GetDB().Begin()
}

func (g *GormConn) RollBack() {
	g.tx.Rollback()
}

func (g *GormConn) Commit() {
	g.tx.Commit()
}
