package dao

import (
	"context"
	"fmt"
	"test.com/project-project/internal/data/project"
	gorms "test.com/project-project/internal/database/gorm"
)

type ProjectDao struct {
	conn *gorms.GormConn
}

func (p ProjectDao) FindMyProjectByMemId(ctx context.Context, memId int64, page int64, size int64, condition string) ([]*project.ProjectAndMenber, int64, error) {
	var pm []*project.ProjectAndMenber
	session := p.conn.Session(ctx)
	index := (page - 1) * size

	sql := fmt.Sprintf("select * from ms_project a,ms_project_member b where a.id=b.project_code and b.member_code=? %s order by sort limit ?,? ", condition)
	raw := session.Raw(sql, memId, index, size)
	raw.Scan(&pm)
	var total int64
	countSql := fmt.Sprintf("select count(*) from ms_project a,ms_project_member b where a.id=b.project_code and b.member_code=? %s ", condition)
	t := session.Raw(countSql, memId)
	err := t.Scan(&total).Error
	return pm, total, err
}

func (p ProjectDao) FindMyCollectedProjectByMemId(ctx context.Context, memId int64, page int64, size int64) ([]*project.ProjectAndMenber, int64, error) {
	var pm []*project.ProjectAndMenber
	session := p.conn.Session(ctx)
	index := (page - 1) * size

	raw := session.Raw("select * from ms_project a,ms_project_collection b where a.id=b.project_code and b.member_code=? order by sort limit ?,? ", memId, index, size)
	raw.Scan(&pm)

	var total int64
	err := session.Model(&project.MemberCollectedProject{}).Where("member_code=?", memId).Count(&total).Error
	return pm, total, err
}

func (p ProjectDao) FindProjectTemplateByCondition(ctx context.Context, page int64, size int64, condition string) ([]*project.ProjectTemplateAll, int64, error) {
	dbResult := []*project.ProjectTemplateAll{}
	session := p.conn.Session(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select a.*,b.name from ms_project_template a,ms_task_stages_template b where a.id=b.project_template_code %s order by a.sort,b.sort limit ?,?", condition)
	err := session.Raw(sql, index, size).Scan(&dbResult).Error

	var total int64
	countSql := fmt.Sprintf("select count(*) from ms_project_template a,ms_task_stages_template b where a.id=b.project_template_code %s order by a.sort,b.sort limit ?,?", condition)
	err = session.Raw(countSql, index, size).Scan(&total).Error
	return dbResult, total, err
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{gorms.New()}
}
