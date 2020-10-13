package public

import (
	"golang.org/x/time/rate"
	"sync"
)

var FlowLimiterHandler *FlowLimiter

type FlowLimiter struct {
	FlowLimiterMap   map[string]*FlowLimiterItem
	FlowLimiterSlice []*FlowLimiterItem
	Locker           sync.RWMutex
}

type FlowLimiterItem struct {
	ServiceName string
	Limiter     *rate.Limiter
}

func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		FlowLimiterMap:   map[string]*FlowLimiterItem{},
		FlowLimiterSlice: []*FlowLimiterItem{},
		Locker:           sync.RWMutex{},
	}
}

// 使用init()可以实现对象单例化
func init() {
	FlowLimiterHandler = NewFlowLimiter()
}

func (limiter *FlowLimiter) GetLimiter(serviceName string, qps float64) (*rate.Limiter, error) {
	for _, item := range limiter.FlowLimiterSlice {
		if item.ServiceName == serviceName {
			return item.Limiter, nil
		}
	}

	// NewLimiter returns a new Limiter that allows events up to rate r and permits bursts of at most b tokens.
	// Limit defines the maximum frequency of some events.
	// Limit is represented as number of events per second.
	// A zero Limit allows no events.
	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	item := &FlowLimiterItem{
		ServiceName: serviceName,
		Limiter:     newLimiter,
	}

	limiter.FlowLimiterSlice = append(limiter.FlowLimiterSlice, item)
	limiter.Locker.Lock()
	defer limiter.Locker.Unlock()
	limiter.FlowLimiterMap[serviceName] = item

	return newLimiter, nil
}
