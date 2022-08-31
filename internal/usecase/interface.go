package usecase

import (
	"context"

	"github.com/cut4cut/hezzl-test-work/internal/entity"
)

type UserRepoInterface interface {
	Create(context.Context, string) (int64, error)
	Delete(context.Context, int64) (int64, error)
	GetList(context.Context, entity.Pagination) ([]entity.User, error)
}

type CacheInterface interface {
	SetValue(context.Context, string, interface{}) error
	GetValue(context.Context, string, interface{}) error
}

type ProducerInterface interface {
	ProduceStruct(interface{}, string) error
}
