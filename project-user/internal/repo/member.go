package repo

import (
	"context"
	"test.com/project-user/internal/data/member"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByAccount(ctx context.Context, account string) (bool, error)
	GetMemberByMobile(ctx context.Context, account string) (bool, error)
	SaveMember(ctx context.Context, member *member.Member) error
}
