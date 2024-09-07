package repo

import (
	"context"
	"test.com/project-project/internal/data/project"
)

type ProjectRepo interface {
	FindMyProjectByMemId(ctx context.Context, menId int64, page int64, size int64) ([]*project.ProjectAndMenber, int64, error)
}
