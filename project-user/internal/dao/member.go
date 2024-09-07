package dao

import (
	"context"
	"gorm.io/gorm"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/database"
	"test.com/project-user/internal/database/gorm"
)

type MemberDao struct { //实现接口
	conn *gorms.GormConn
}

func NewMemberDao() *MemberDao { //新建member数据库操作类
	return &MemberDao{
		gorms.New(),
	}
}

func (m *MemberDao) GetMemberByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("email=?", email).Count(&count).Error
	//model()找到哪个表
	return count > 0, err
}

func (m *MemberDao) GetMemberByAccount(ctx context.Context, account string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("account=?", account).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) GetMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("mobile=?", mobile).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) SaveMember(conn database.DbConn, ctx context.Context, member *member.Member) error { //增加
	m.conn = conn.(*gorms.GormConn) //因为是数据库连接接口的具体实现，所以将接口转为gorm的数据库连接，向下转型

	return m.conn.Tx(ctx).Create(member).Error
}

func (m *MemberDao) FindMember(ctx context.Context, account string, password string) (*member.Member, error) {
	var mem *member.Member
	err := m.conn.Session(ctx).Where("account=? and password=?", account, password).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return mem, err
}

func (m *MemberDao) FindMemberById(ctx context.Context, id int64) (mem *member.Member, err error) { //注意返回值还有变量名的写法
	err = m.conn.Session(ctx).Where("id=?", id).First(&mem).Error
	return
}
