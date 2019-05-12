package splash

import (
	"sync"
)

type Pool struct {
	size   uint
	queue  chan job
	errors chan error
	kill   chan interface{}
	wg     sync.WaitGroup
}

func (pool *Pool) enqueue(j job) {
	pool.queue <- j
}

func (pool *Pool) killOne() {
	pool.kill <- struct{}{}
	pool.size--
}

func (pool *Pool) killSome(amount int) {
	for i := 0; i < amount; i++ {
		pool.killOne()
	}
}

func (pool *Pool) killAll() {
	pool.killSome(int(pool.size))
}

func (pool *Pool) rawAlloc() {
	newWorker(pool.queue, pool.errors, pool.kill, &pool.wg)
}

func (pool *Pool) allocOne() {
	pool.rawAlloc()
	pool.size++
}

func (pool *Pool) allocSome(amount int) {
	for i := 0; i < amount; i++ {
		pool.rawAlloc()
	}
	pool.size += uint(amount)
}
