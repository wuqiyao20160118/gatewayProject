package public

import (
	"sync"
	"time"
)

var FlowCounterHandler *FlowCounter

type FlowCounter struct {
	RedisFlowCountMap   map[string]*RedisFlowCountService
	RedisFlowCountSlice []*RedisFlowCountService
	Locker              sync.RWMutex
}

func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap:   map[string]*RedisFlowCountService{},
		RedisFlowCountSlice: []*RedisFlowCountService{},
		Locker:              sync.RWMutex{},
	}
}

// 使用init()可以实现对象单例化
func init() {
	FlowCounterHandler = NewFlowCounter()
}

func (counter *FlowCounter) GetFlowCounter(serviceName string) (*RedisFlowCountService, error) {
	for _, item := range counter.RedisFlowCountSlice {
		if item.AppID == serviceName {
			return item, nil
		}
	}

	newCounter := NewRedisFlowCountService(serviceName, time.Second)

	counter.RedisFlowCountSlice = append(counter.RedisFlowCountSlice, newCounter)
	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.RedisFlowCountMap[serviceName] = newCounter

	return newCounter, nil
}
