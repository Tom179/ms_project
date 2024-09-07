package dao

import (
	"context"
	"test.com/project-project/internal/data/menu"
	gorms "test.com/project-project/internal/database/gorm"
)

type MenuDao struct {
	conn *gorms.GormConn
}

func (m *MenuDao) FindMenus(ctx context.Context) (pms []*menu.ProjectMenu, err error) {
	err = m.conn.Session(ctx).Order("pid,sort asc,id asc").Find(&pms).Error
	return
}

func NewMenuDao() *MenuDao {
	return &MenuDao{
		gorms.New(),
	}
}
