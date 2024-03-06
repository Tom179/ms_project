package project_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	project "test.com/project-grpc/project"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data/menu"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/repo"
)

type ProjectService struct {
	project.UnimplementedProjectServiceServer
	cache       repo.Cache
	menuRepo    repo.MenuRepo
	transaction tran.Transaction
}

func New() *ProjectService {
	return &ProjectService{
		cache:       dao.Rc,
		menuRepo:    dao.NewMenuDao(),
		transaction: dao.NewTransaction(),
	}
}

func (p *ProjectService) Index(context.Context, *project.IndexMessage) (*project.IndexResponse, error) {
	pms, err := p.menuRepo.FindMenus(context.Background())
	if err != nil {
		zap.L().Error("Index db find Menu error", zap.Error(err))
		return nil, err
	}

	projectTree := menu.CovertChild(pms)
	var mms []*project.MenuMessage //grpc交互信息
	copier.Copy(mms, projectTree)
	return &project.IndexResponse{Menus: mms}, nil
}
