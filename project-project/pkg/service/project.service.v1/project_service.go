package project_service_v1

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"strconv"
	common "test.com/project-common"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	project_grpc "test.com/project-grpc/project"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data/menu"
	"test.com/project-project/internal/data/project"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/repo"
	"test.com/project-project/pkg/model"
	"time"
)

type ProjectService struct {
	project_grpc.UnimplementedProjectServiceServer
	cache       repo.Cache
	menuRepo    repo.MenuRepo
	projectRepo repo.ProjectRepo
	transaction tran.Transaction
}

func (p *ProjectService) ReadProject(ctx context.Context, msg *project_grpc.ProjectMessage) (*project_grpc.ProjectMessage, error) {
	projectId := msg.ProjectCode
	rpcRsp := project_grpc.ProjectMessage{}
	singleProject, err := p.projectRepo.ReadOneProject(ctx, projectId)
	if err != nil {
		return nil, errs.GrpcError(err.(*errs.BError))
	}
	copier.Copy(&rpcRsp, singleProject)
	return &rpcRsp, nil
}

func (p *ProjectService) CreateProject(ctx context.Context, msg *project_grpc.ProjectMessage) (*project_grpc.ProjectMessage, error) {
	fmt.Println("创建项目的用户id为:", msg.Id)
	project := project.Project{
		Name:              msg.Name,
		TemplateCode:      msg.TemplateCode,
		Description:       msg.Description,
		OrganizationCode:  msg.OrganizationCode,
		Cover:             "https://images.pexels.com/photos/372787/pexels-photo-372787.jpeg",
		CreateTime:        time.Now().UnixMilli(),
		AccessControlType: model.Open,
		Archive:           model.NoArcheve,
		TaskBoardTheme:    model.Simple,
	}

	err := p.transaction.Action(func(conn database.DbConn) error {
		err := p.projectRepo.CreateProject(ctx, &project, msg.Id)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, errs.GrpcError(err.(*errs.BError)) //????会报错类型转换错误
	} //panic: interface conversion: error is *mysql.MySQLError, not *errs.BError

	rpcRsp := project_grpc.ProjectMessage{}
	copier.Copy(&rpcRsp, project)
	rpcRsp.CreateTime = strconv.FormatInt(project.CreateTime, 10)
	return &rpcRsp, nil
}

func New() *ProjectService {
	return &ProjectService{
		cache:       dao.Rc,
		menuRepo:    dao.NewMenuDao(),
		projectRepo: dao.NewProjectDao(),
		transaction: dao.NewTransaction(),
	}
}

func (p *ProjectService) Index(context.Context, *project_grpc.IndexMessage) (*project_grpc.IndexResponse, error) {
	pms, err := p.menuRepo.FindMenus(context.Background())
	if err != nil {
		zap.L().Error("Index db find Menu error", zap.Error(err))
		return nil, err
	}

	projectTree := menu.GetChild(pms)
	var mms []*project_grpc.MenuMessage //grpc交互信息
	copier.Copy(&mms, projectTree)
	return &project_grpc.IndexResponse{Menus: mms}, nil
}

func (p *ProjectService) FindProjectByMemId(ctx context.Context, msg *project_grpc.ProjectRpcMessage) (*project_grpc.MyProjectResponse, error) {
	memId := msg.MemberId
	page := msg.Page
	pageSize := msg.PageSize
	var pms []*project.ProjectAndMenber
	var total int64
	var err error
	if msg.SelectBy == "" || msg.SelectBy == "my" {
		pms, total, err = p.projectRepo.FindMyProjectByMemId(context.Background(), memId, page, pageSize, "")
	}
	if msg.SelectBy == "archive" {
		//因为这里的输入不是用户前端传的，字符串直接拼接sql是我们自己干的，所以不存在注入攻击
		pms, total, err = p.projectRepo.FindMyProjectByMemId(context.Background(), memId, page, pageSize, "and archive=1")

	}
	if msg.SelectBy == "deleted" {
		pms, total, err = p.projectRepo.FindMyProjectByMemId(context.Background(), memId, page, pageSize, "and deleted=1")

	}
	if msg.SelectBy == "collect" {
		pms, total, err = p.projectRepo.FindMyCollectedProjectByMemId(context.Background(), memId, page, pageSize)
	}

	if err != nil {
		zap.L().Error("menu findMyProject error:", zap.Error(err))
		fmt.Println("myProject数据库查询失败")
		return nil, err
	}
	if pms == nil {
		return &project_grpc.MyProjectResponse{Pm: []*project_grpc.ProjectMessage{}, Total: total}, nil
	}

	var rsp []*project_grpc.ProjectMessage
	copier.Copy(&rsp, pms)
	for _, v := range rsp {

		v.Code, _ = encrypts.EncryptInt64(v.Id, encrypts.AESKEY)
		pam := project.ToMap(pms)[v.Id] //根据id得到ProjectAndMember
		organizationCode, _ := encrypts.EncryptInt64(pam.OrganizationCode, encrypts.AESKEY)
		organizationCodeInt64, _ := strconv.ParseInt(organizationCode, 10, 64)

		v.AccessControlType = pam.GetAccessControlType()
		v.OrganizationCode = organizationCodeInt64
		v.JoinTime = common.FormatByMill(pam.JoinTime)
		v.OwnerName = msg.MemberName
		v.Order = int32(pam.Sort)
		v.CreateTime = common.FormatByMill(pam.JoinTime)
	}
	return &project_grpc.MyProjectResponse{Pm: rsp, Total: total}, nil
}

func (p *ProjectService) FindProjectTemplate(ctx context.Context, msg *project_grpc.ProjectRpcMessage) (*project_grpc.ProjectTemplateResponse, error) {
	rsp := project_grpc.ProjectTemplateResponse{}

	var condition string
	var dbResult []*project.ProjectTemplateAll
	var total int64
	var err error

	if msg.ViewType == 0 { //自定义模板
		condition = fmt.Sprintf("where member_code=%d", int(msg.MemberId)) //和底层DB耦合，业务层还需要弄清表结构，不太好
	} else if msg.ViewType == 1 { //系统模板
		condition = fmt.Sprintf("where is_system=1")
	} else if msg.ViewType == -1 { //所有模板
		condition = fmt.Sprintf("")
	}

	dbResult, total, err = p.projectRepo.FindProjectTemplateByCondition(ctx, msg.Page, msg.PageSize, condition)

	err = copier.Copy(&rsp.Ptm, dbResult)
	if err != nil {
		fmt.Println("拷贝失败：", err)
	}
	for _, t := range rsp.Ptm {
		t.Code, err = encrypts.Encrypt(string(t.Id), encrypts.AESKEY)
	}
	rsp.Total = total
	fmt.Println("rpcServce最终构建的rsp为", rsp.Ptm)
	fmt.Println("dbREsult", dbResult)

	return &rsp, err //rpc服务方
}
