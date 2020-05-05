package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/teploff/antibruteforce/internal/implementation/repository"
)

const (
	absentKey  = "absent_key"
	bucketKey1 = "bucket_key_1"
	bucketKey2 = "bucket_key_2"
)

func TestGetAbsentBucketFromEmptyBasket(t *testing.T) {
	rate := 1
	expireTime := time.Second
	basket := repository.NewLeakyBucket(rate, expireTime)

	limiter, err := basket.Get(absentKey)
	assert.Error(t, err)
	assert.Nil(t, limiter)
}

func TestGetAbsentBucketFromNotEmptyBasket(t *testing.T) {
	rate := 1
	expireTime := time.Second
	basket := repository.NewLeakyBucket(rate, expireTime)

	limiter, err := basket.Add(bucketKey1)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	limiter, err = basket.Get(absentKey)
	assert.Error(t, err)
	assert.Nil(t, limiter)
}

func TestGetExistBucket(t *testing.T) {
	rate := 1
	expireTime := time.Second
	basket := repository.NewLeakyBucket(rate, expireTime)

	limiter, err := basket.Add(bucketKey1)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	limiter, err = basket.Get(bucketKey1)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)
}

func TestAddExistBucket(t *testing.T) {
	rate := 1
	expireTime := time.Second
	basket := repository.NewLeakyBucket(rate, expireTime)

	limiter, err := basket.Add(bucketKey1)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	limiter, err = basket.Add(bucketKey1)
	assert.Error(t, err)
	assert.Nil(t, limiter)
}

func TestDeleteAbsentBucketFromEmptyBasket(t *testing.T) {
	rate := 1
	expireTime := time.Second
	basket := repository.NewLeakyBucket(rate, expireTime)

	err := basket.Delete(absentKey)
	assert.Error(t, err)
}

func TestDeleteAbsentBucketFromNotEmptyBasket(t *testing.T) {
	rate := 1
	expireTime := time.Second
	basket := repository.NewLeakyBucket(rate, expireTime)

	limiter, err := basket.Add(bucketKey1)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	err = basket.Delete(absentKey)
	assert.Error(t, err)
}

func TestDeleteExistBucket(t *testing.T) {
	rate := 1
	expireTime := time.Second
	basket := repository.NewLeakyBucket(rate, expireTime)

	limiter, err := basket.Add(bucketKey1)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	err = basket.Delete(bucketKey1)
	assert.NoError(t, err)

	limiter, err = basket.Get(bucketKey1)
	assert.Nil(t, limiter)
	assert.Error(t, err)
}

func TestNoOneBucketsShouldClean(t *testing.T) {
	rate := 1
	expireTime := time.Millisecond * 100
	basket := repository.NewLeakyBucket(rate, expireTime)
	ticker := time.NewTicker(time.Millisecond * 10)

	limiter, err := basket.Add(bucketKey1)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	limiter, err = basket.Add(bucketKey2)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	<-ticker.C
	basket.Clean()

	limiter, err = basket.Get(bucketKey1)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	limiter, err = basket.Get(bucketKey2)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)
}

func TestOneBucketsShouldClean(t *testing.T) {
	rate := 1
	expireTime := time.Millisecond * 50
	basket := repository.NewLeakyBucket(rate, expireTime)
	ticker := time.NewTicker(time.Millisecond * 100)

	limiter, err := basket.Add(bucketKey1)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	<-ticker.C

	limiter, err = basket.Add(bucketKey2)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	basket.Clean()

	limiter, err = basket.Get(bucketKey1)
	assert.Error(t, err)
	assert.Nil(t, limiter)

	limiter, err = basket.Get(bucketKey2)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)
}

func TestAllBucketsShouldClean(t *testing.T) {
	rate := 1
	expireTime := time.Millisecond * 50
	basket := repository.NewLeakyBucket(rate, expireTime)
	ticker := time.NewTicker(time.Millisecond * 100)

	limiter, err := basket.Add(bucketKey1)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	limiter, err = basket.Add(bucketKey2)
	assert.NoError(t, err)
	assert.NotNil(t, limiter)

	<-ticker.C

	basket.Clean()

	limiter, err = basket.Get(bucketKey1)
	assert.Error(t, err)
	assert.Nil(t, limiter)

	limiter, err = basket.Get(bucketKey2)
	assert.Error(t, err)
	assert.Nil(t, limiter)
}
