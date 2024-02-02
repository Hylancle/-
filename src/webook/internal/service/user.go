// 业务代码
package service

import (
	"context"
	"xiaoweishu/webook/internal/domain"
	"xiaoweishu/webook/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
func (svc *UserService) Signup(ctx context.Context, user domain.User) error {
	return svc.repo.Create(ctx, user)
}
