package splash

import "sync"

func NewPool(poolSize uint, bufferSize uint) *Pool {
	queue := make(chan job, bufferSize)
	errors := make(chan error, 10)
	kill := make(chan interface{}, 10)
	wg := sync.WaitGroup{}

	pool := Pool{
		size:   poolSize,
		queue:  queue,
		errors: errors,
		kill:   kill,
		wg:     wg,
	}

	defer func() {
		pool.allocSome(int(poolSize))
	}()

	return &pool
}

func (pool *Pool) Grow(delta uint) {
	pool.wg.Add(int(delta))

	var i uint
	for i = 0; i < delta; i++ {
		go func() {
			defer pool.wg.Done()

			pool.allocOne()
		}()
	}
}

func (pool *Pool) Shrink(delta uint) {
	if delta > pool.size {
		pool.killAll()
	} else {
		pool.killSome(int(delta))
	}
}

func (pool *Pool) Exec(f func(...interface{}) error, args ...interface{}) {
	pool.wg.Add(1)
	go func() {
		defer pool.wg.Done()

		pool.enqueue(newJob(f, args))
	}()
}

func (pool *Pool) ExecNil(f func() error) {
	pool.Exec(func(_ ...interface{}) error {
		return f()
	}, nil)
}

func (pool *Pool) SyncExec(f func(...interface{}) error, args ...interface{}) {
	pool.enqueue(newJob(f, args))
}

func (pool *Pool) SyncExecNil(f func() error) {
	pool.SyncExec(func(_ ...interface{}) error {
		return f()
	}, nil)
}
