package balance

import (
	"errors"
	"math/rand"
	"sync/atomic"
)

var ErrNoPointer = errors.New("no instance available")

type LoadBalance interface {
	Target() (string, error)
}

type RandomBalance struct {
	instances []string
	rand *rand.Rand
	lens int
}

type RoundRobinBalance struct {
	instances []string
	current *uint64
}

func (r *RandomBalance) Target() (string, error) {
	if len(r.instances) <= 0 {
		return "", ErrNoPointer
	}
	return r.instances[r.rand.Intn(r.lens)], nil
}

func (rb *RoundRobinBalance) Target() (string, error) {
	lens := len(rb.instances)
	if lens <= 0 {
		return "", ErrNoPointer
	}
	old := atomic.AddUint64(rb.current, 1) - 1
	idx := old % uint64(len(rb.instances))
	return rb.instances[idx], nil
}

//随机
func NewRandom(s []string, seed int64) LoadBalance {
	return &RandomBalance{
		instances: s,
		rand: rand.New(rand.NewSource(seed)),
		lens: len(s),
	}
}

//轮训
func NewRoundRobin(s []string) LoadBalance {
	return &RoundRobinBalance{
		instances: s,
		current: new(uint64),
	}
}