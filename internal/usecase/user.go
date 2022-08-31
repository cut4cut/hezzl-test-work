package usecase

import (
	"context"
	"fmt"

	"github.com/cut4cut/hezzl-test-work/internal/entity"
	"github.com/go-redis/redis/v8"
)

type UserUseCase struct {
	repo     UserRepoInterface
	cache    CacheInterface
	producer ProducerInterface
}

func New(r UserRepoInterface, c CacheInterface, p ProducerInterface) *UserUseCase {
	return &UserUseCase{
		repo:     r,
		cache:    c,
		producer: p,
	}
}

func (uc *UserUseCase) Create(ctx context.Context, name string) (int64, error) {
	var (
		id   int64 = -1
		user entity.User
		err  error
	)

	if user, err = entity.NewUser(name); err != nil {
		return id, fmt.Errorf("userUseCase - create - entity.NewUser: %w", err)
	}

	if id, err = uc.repo.Create(ctx, user.Name); err != nil {
		return id, fmt.Errorf("userUseCase - create - uc.repo.Create: %w", err)
	}

	user.Id = id
	if err = uc.producer.ProduceStruct(user, "user"); err != nil {
		return id, fmt.Errorf("userUseCase - create - uc.producer.ProduceStruct: %w", err)
	}

	return id, nil
}

func (uc *UserUseCase) Delete(ctx context.Context, id int64) (int64, error) {

	if id, err := uc.repo.Delete(ctx, id); err != nil {
		return id, fmt.Errorf("userUseCase - delete - uc.repo.Dreate: %w", err)
	}

	return id, nil
}

func (uc *UserUseCase) GetList(ctx context.Context, pag entity.Pagination) ([]entity.User, error) {
	var (
		users []entity.User
		err   error
	)

	if err = uc.cache.GetValue(ctx, pag.Key, &users); err == redis.Nil {
		users, err = uc.repo.GetList(ctx, pag)
		if err != nil {
			return users, fmt.Errorf("userUseCase - getList - uc.repo.GetList: %w", err)
		}

		if err = uc.cache.SetValue(ctx, pag.Key, &users); err != nil {
			return users, fmt.Errorf("userUseCase - getList - uc.cache.SetValue: %w", err)
		}

	} else if err != nil {
		return users, fmt.Errorf("userUseCase - getList - uc.cache.GetValue: %w", err)
	}

	return users, nil
}
