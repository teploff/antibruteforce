package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/teploff/antibruteforce/internal/domain/entity"
)

func TestLessRequestsPerMinuteThanAllowed(t *testing.T) {
	limiter := entity.NewLimiter(10, time.Second)
	ticker := time.NewTicker(110 * time.Millisecond)
	tickerStop := time.NewTicker(time.Second)
	allowedRequests := make([]bool, 0, 9)
	denyRequests := make([]bool, 0)
LOOP:
	for {
		select {
		case <-ticker.C:
			if limiter.Allow() {
				allowedRequests = append(allowedRequests, true)
			} else {
				denyRequests = append(denyRequests, false)
			}
		case <-tickerStop.C:
			break LOOP
		}
	}
	assert.Equal(t, 9, len(allowedRequests))
	assert.Equal(t, 0, len(denyRequests))
}

func TestMoreRequestsPerMinuteThanAllowed(t *testing.T) {
	limiter := entity.NewLimiter(10, time.Second)
	ticker := time.NewTicker(90 * time.Millisecond)
	tickerStop := time.NewTicker(time.Second)
	allowedRequests := make([]bool, 0, 10)
	denyRequests := make([]bool, 0, 1)
LOOP:
	for {
		select {
		case <-ticker.C:
			if limiter.Allow() {
				allowedRequests = append(allowedRequests, true)
			} else {
				denyRequests = append(denyRequests, false)
			}
		case <-tickerStop.C:
			break LOOP
		}
	}
	assert.Equal(t, 10, len(allowedRequests))
	assert.Equal(t, 1, len(denyRequests))
}
