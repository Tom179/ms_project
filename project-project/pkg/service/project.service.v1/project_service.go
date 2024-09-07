package project_service_v1

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	common "test.com/project-common"
	"test.com/project-common/encrypts"
	project_grpc "test.com/project-grpc/project"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data/menu"
	"test.com/project-project/internal/data/project"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/repo"
)

type ProjectService struct {
	project_grpc.UnimplementedProjectServiceServer
	cache       repo.Cache
	menuRepo    repo.MenuRepo
	projectRepo repo.ProjectRepo
	transaction tran.Transaction
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
	pms, total, err := p.projectRepo.FindMyProjectByMemId(context.Background(), memId, page, pageSize)
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

		v.AccessControlType = pam.GetAccessControlType()
		v.OrganizationCode, _ = encrypts.EncryptInt64(pam.OrganizationCode, encrypts.AESKEY)
		v.JoinTime = common.FormatByMill(pam.JoinTime)
		v.OwnerName = msg.MemberName
		v.Order = int32(pam.Sort)
		v.CreateTime = common.FormatByMill(pam.JoinTime)
	}
	return &project_grpc.MyProjectResponse{Pm: rsp, Total: total}, nil
}
