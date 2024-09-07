package dao

import (
	"context"
	"test.com/project-project/internal/data/project"
	gorms "test.com/project-project/internal/database/gorm"
)

type ProjectDao struct {
	conn *gorms.GormConn
}

func (p ProjectDao) FindMyProjectByMemId(ctx context.Context, memId int64, page int64, size int64) ([]*project.ProjectAndMenber, int64, error) {
	var pm []*project.ProjectAndMenber
	session := p.conn.Session(ctx)
	index := (page - 1) * size
	raw := session.Raw("select * from ms_project a,ms_project_member b where a.id=b.project_code and b.member_code=? order by sort limit ?,? ", memId, index, size)
	raw.Scan(&pm)
	var total int64
	err := session.Model(&project.MemberProject{}).Where("member_code=?", memId).Count(&total).Error
	return pm, total, err
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{gorms.New()}
}
