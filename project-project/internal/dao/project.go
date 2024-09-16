package dao

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"test.com/project-project/internal/data/project"
	"test.com/project-project/internal/data/task"
	gorms "test.com/project-project/internal/database/gorm"
)

type ProjectDao struct {
	conn *gorms.GormConn
}

func (p ProjectDao) ReadOneProject(ctx context.Context, projectID int64) (project.SingleProjectMessage, error) {

	sql := "select a.*,b.member_code,c.name as owner_name,c.avatar as owner_avatar from ms_project a,ms_project_member b,ms_member c where a.id=b.project_code and b.member_code=c.id and project_code=?"
	singleProjectMessage := project.SingleProjectMessage{}
	sql2 := "select count(*) as collected from ms_project_collection where  project_code=? and member_code=?"

	err := p.conn.Session(ctx).Raw(sql, projectID).Find(&singleProjectMessage).Error
	err = p.conn.Session(ctx).Model(&project.MemberCollectedProject{}).Raw(sql2, projectID, singleProjectMessage.MemberCode).Find(&singleProjectMessage).Error
	return singleProjectMessage, err
}

func (p ProjectDao) CreateProject(ctx context.Context, newProject *project.Project, id int64) error {
	err := p.conn.Session(ctx).Create(&newProject).Error
	ralation := project.MemberProject{
		ProjectCode: newProject.Id,
		MemberCode:  id,
		JoinTime:    newProject.CreateTime,
	}
	err = p.conn.Session(ctx).Create(&ralation).Error
	if err != nil {
		return err
	}

	return nil
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
	pt := []*project.ProjectTemplate{}
	session := p.conn.Session(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from ms_project_template %s order by sort limit ?,?", condition)
	err := session.Raw(sql, index, size).Scan(&pt).Error
	copier.Copy(&dbResult, pt)

	//fmt.Println("底层db结果", pt)

	/*	fmt.Println("pt", *pt[0])
		fmt.Println("dbResult", *dbResult[0])
	*/
	for i, templateDB := range pt {
		tt := []task.MsTaskStagesTemplate{}
		err = session.Where("project_template_code=?", templateDB.Id).Find(&tt).Error
		IdTasks := task.CovertProjectMap(tt)
		dbResult[i].TaskStages = IdTasks[dbResult[i].Id]
		dbResult[i] = templateDB.Convert(dbResult[i].TaskStages)
		//将一个ProjectTemplateAll填充完整，转换各种格式
		//fmt.Println(dbResult[i])
	}

	var total int64
	countSql := fmt.Sprintf("select count(*) from ms_project_template %s", condition)
	err = session.Raw(countSql).Scan(&total).Error
	return dbResult, total, err
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{gorms.New()}
}
