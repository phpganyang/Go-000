package roll

import (
	"fmt"
	"sync"
	"time"
)

//借鉴hystrix-go思想实现,默认10个buckets
type RollingDemo struct {
	Buckets map[int64]*bucket
	Mu      *sync.RWMutex
}

//构建计数桶
type bucket struct {
	Value int64
}

//初始化桶
func NewBucket() *RollingDemo {
	return &RollingDemo{
		Buckets: make(map[int64]*bucket),
		Mu:      &sync.RWMutex{},
	}
}

//获取当前桶处于第几个
func (rd *RollingDemo) getNowBucket() *bucket {
	now := time.Now().Unix()
	var b *bucket
	var ok bool
	if b, ok = rd.Buckets[now]; !ok {
		b = &bucket{}
		rd.Buckets[now] = b
	}
	return b
}

func (rd *RollingDemo) removeOldBuckets() {
	expired := time.Now().Unix() - 10
	for timestamp := range rd.Buckets {
		//fmt.Print(rd.Buckets)
		if timestamp <= expired {
			delete(rd.Buckets, timestamp)
		}
	}
}

// Increment 累加最新桶的计数器
func (rd *RollingDemo) Increment() {
	rd.Mu.Lock()
	b := rd.getNowBucket()
	b.Value++
	fmt.Print(b)
	rd.removeOldBuckets()

	rd.Mu.Unlock()
}

// UpdateMax 将最新桶的计数器置为某个最大值
func (rd *RollingDemo) UpdateMax(n int64) {
	rd.Mu.Lock()
	b := rd.getNowBucket()
	if n > b.Value {
		b.Value = n
	}
	rd.Mu.Unlock()
	rd.removeOldBuckets()
}

// Sum 计算最新 10 个桶内计数器的和
func (rd *RollingDemo) Sum(now time.Time) int64 {
	sum := int64(0)

	rd.Mu.RLock()
	defer rd.Mu.RUnlock()

	for timestamp, bucket := range rd.Buckets {
		if timestamp >= now.Unix()-10 {
			sum += bucket.Value
		}
	}
	return sum
}

// Max 获取最新 10 个桶内计数器的最大值
func (rd *RollingDemo) Max(now time.Time) int64 {
	var max int64

	rd.Mu.RLock()
	defer rd.Mu.RUnlock()

	for timestamp, bucket := range rd.Buckets {
		if timestamp >= now.Unix()-10 {
			if bucket.Value > max {
				max = bucket.Value
			}
		}
	}
	return max
}

// Avg 计算最新 10 个桶内计数器的平均值
func (rd *RollingDemo) Avg(now time.Time) int64 {
	return rd.Sum(now) / 10
}
