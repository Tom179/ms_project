package repo

import (
	"context"
	"test.com/project-project/internal/data/project"
)

type ProjectRepo interface {
	FindMyProjectByMemId(ctx context.Context, menId int64, page int64, size int64, condition string) ([]*project.ProjectAndMenber, int64, error)
	FindMyCollectedProjectByMemId(ctx context.Context, memId int64, page int64, size int64) ([]*project.ProjectAndMenber, int64, error)
	FindProjectTemplateByCondition(ctx context.Context, page int64, size int64, condition string) ([]*project.ProjectTemplateAll, int64, error)
}
