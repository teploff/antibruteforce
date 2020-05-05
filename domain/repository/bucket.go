package repository

import (
	"github.com/teploff/antibruteforce/domain/entity"
)

type BucketStorable interface {
	Add(key string) (*entity.Limiter, error)
	Delete(key string) error
	Get(key string) (*entity.Limiter, error)
	Clean()
}
