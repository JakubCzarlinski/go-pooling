package pooling

import (
	"bytes"
	"sync"
)

type Poolable[T any, A any] interface {
	Reset(args A)
}

type Pool[T Poolable[T, A], A any] struct {
	pool sync.Pool
}

func NewPool[T Poolable[T, A], A any](new func() T) *Pool[T, A] {
	return &Pool[T, A]{
		pool: sync.Pool{
			New: func() interface{} {
				return new()
			},
		},
	}
}

func (p *Pool[T, A]) Get() T {
	return p.pool.Get().(T)
}

func (p *Pool[T, A]) Put(obj T) {
	p.pool.Put(obj)
}

func (p *Pool[T, A]) Reset(obj T, args A) {
	obj.Reset(args)
	p.pool.Put(obj)
}

type BytesBuffer struct {
	*bytes.Buffer
}

func (b BytesBuffer) Reset(struct{}) {
	b.Buffer.Reset()
}

func CreateBytesBufferPool(size int) *Pool[*BytesBuffer, struct{}] {
	return NewPool(func() *BytesBuffer {
		return &BytesBuffer{bytes.NewBuffer(make([]byte, 0, size))}
	})
}
