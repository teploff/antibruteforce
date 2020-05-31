package middleware

import (
	"fmt"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/teploff/antibruteforce/internal/domain/service"
	"net"
	"time"
)

type prometheusMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           service.AdminService
}

func NewPrometheusMiddleware(next service.AdminService) service.AdminService {
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "antibruteforce",
		Subsystem: "admin_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "antibruteforce",
		Subsystem: "admin_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "antibruteforce",
		Subsystem: "admin_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})
	return &prometheusMiddleware{requestCount: requestCount, requestLatency: requestLatency, countResult: countResult, next: next}
}

func (p prometheusMiddleware) ResetBucketByLogin(login string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "reset_bucket_by_login", "error", fmt.Sprint(err != nil)}
		p.requestCount.With(lvs...).Add(1)
		p.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = p.next.ResetBucketByLogin(login)
	return
}

func (p prometheusMiddleware) ResetBucketByPassword(password string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "reset_bucket_by_password", "error", fmt.Sprint(err != nil)}
		p.requestCount.With(lvs...).Add(1)
		p.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = p.next.ResetBucketByPassword(password)
	return
}

func (p prometheusMiddleware) ResetBucketByIP(ip net.IP) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "reset_bucket_by_ip", "error", fmt.Sprint(err != nil)}
		p.requestCount.With(lvs...).Add(1)
		p.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = p.next.ResetBucketByIP(ip)
	return
}

func (p prometheusMiddleware) AddInBlacklist(ipNet *net.IPNet) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "add_in_blacklist", "error", fmt.Sprint(err != nil)}
		p.requestCount.With(lvs...).Add(1)
		p.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = p.next.AddInBlacklist(ipNet)
	return
}

func (p prometheusMiddleware) RemoveFromBlacklist(ipNet *net.IPNet) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "remove_from_blacklist", "error", fmt.Sprint(err != nil)}
		p.requestCount.With(lvs...).Add(1)
		p.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = p.next.RemoveFromBlacklist(ipNet)
	return
}

func (p prometheusMiddleware) AddInWhitelist(ipNet *net.IPNet) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "add_in_whitelist", "error", fmt.Sprint(err != nil)}
		p.requestCount.With(lvs...).Add(1)
		p.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = p.next.AddInWhitelist(ipNet)
	return
}

func (p prometheusMiddleware) RemoveFromWhitelist(ipNet *net.IPNet) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "remove_from_whitelist", "error", fmt.Sprint(err != nil)}
		p.requestCount.With(lvs...).Add(1)
		p.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = p.next.RemoveFromWhitelist(ipNet)
	return
}
