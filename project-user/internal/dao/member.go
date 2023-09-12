package dao

import (
	"context"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/database/gorm"
)

type MemberDao struct { //实现接口
	conn *gorm.GormConn
}

func NewMemberDao() *MemberDao { //新建member数据库操作类
	return &MemberDao{
		gorm.New(),
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

func (m *MemberDao) SaveMember(ctx context.Context, member *member.Member) error { //增加
	return m.conn.Session(ctx).Create(member).Error
}
