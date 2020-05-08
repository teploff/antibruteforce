package repository

import (
	"github.com/teploff/antibruteforce/domain/entity"
)

// BucketStorable provides storable interface for leaky buckets
//
// Add - adding bucket by key
//
// Delete - deleting bucket by key
//
// Get - getting bucket by key
//
// Clean - flush bucket.
type BucketStorable interface {
	Add(key string) (*entity.Limiter, error)
	Delete(key string) error
	Get(key string) (*entity.Limiter, error)
	Clean()
}
