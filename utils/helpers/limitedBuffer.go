package helpers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/snappy"
	"github.com/valyala/bytebufferpool"
	"golang.org/x/sync/semaphore"
	"io"
	"strconv"
)

type RateLimitedBuffer interface {
	Bytes() []byte
	Write([]byte) (int, error)
	Release()
}

type RateLimitedPooledBuffer struct {
	pool   *RateLimitedPool
	limit  int
	buffer *bytebufferpool.ByteBuffer
}

func (l *RateLimitedPooledBuffer) Write(msg []byte) (int, error) {
	if len(msg)+l.buffer.Len() > l.limit {
		return 0, fmt.Errorf("buffer size overflow")
	}
	return l.buffer.Write(msg)
}

func (l *RateLimitedPooledBuffer) Bytes() []byte {
	return l.buffer.Bytes()
}

func (l *RateLimitedPooledBuffer) Release() {
	l.pool.releasePooledBuffer(l)
}

type RateLimitedSliceBuffer struct {
	bytes []byte
	pool  *RateLimitedPool
}

func (l *RateLimitedSliceBuffer) Write(msg []byte) (int, error) {
	return 0, nil
}

func (l *RateLimitedSliceBuffer) Bytes() []byte {
	return l.bytes
}

func (l *RateLimitedSliceBuffer) Release() {
	l.pool.releaseSlice(l)
}

type RateLimitedPool struct {
	limit       int
	rateLimiter *semaphore.Weighted
}

func (r *RateLimitedPool) acquirePooledBuffer(limit int) (RateLimitedBuffer, error) {
	if limit > r.limit {
		return nil, fmt.Errorf("limit too big")
	}
	err := r.rateLimiter.Acquire(context.Background(), int64(limit))
	if err != nil {
		return nil, err
	}
	return &RateLimitedPooledBuffer{
		pool:   r,
		limit:  limit,
		buffer: bytebufferpool.Get(),
	}, nil
}

func (r *RateLimitedPool) releasePooledBuffer(buffer *RateLimitedPooledBuffer) {
	r.rateLimiter.Release(int64(buffer.limit))
	bytebufferpool.Put(buffer.buffer)
}

func (r *RateLimitedPool) acquireSlice(size int) (RateLimitedBuffer, error) {
	if size > r.limit {
		return nil, fmt.Errorf("size too big")
	}
	err := r.rateLimiter.Acquire(context.Background(), int64(size))
	if err != nil {
		return nil, err
	}
	return &RateLimitedSliceBuffer{
		bytes: make([]byte, size),
		pool:  r,
	}, nil
}

func (r *RateLimitedPool) releaseSlice(buffer *RateLimitedSliceBuffer) {
	r.rateLimiter.Release(int64(len(buffer.bytes)))
}

var requestPool = RateLimitedPool{
	limit:       50 * 1024 * 1024,
	rateLimiter: semaphore.NewWeighted(50 * 1024 * 1024),
}
var pbPool = RateLimitedPool{
	limit:       50 * 1024 * 1024,
	rateLimiter: semaphore.NewWeighted(50 * 1024 * 1024),
}

func GetRawBody(ctx *fiber.Ctx) (RateLimitedBuffer, error) {
	if ctx.Get("content-length", "") == "" {
		return nil, fmt.Errorf("content-length is required")
	}
	ctxLen, _ := strconv.Atoi(ctx.Get("content-length", ""))
	buf, err := requestPool.acquirePooledBuffer(ctxLen)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(buf, ctx.Context().RequestBodyStream())
	if err != nil {
		return nil, err
	}
	if ctx.Get("content-type", "") != "application/x-protobuf" {
		return buf, err
	}
	defer buf.Release()
	decompSize, err := snappy.DecodedLen(buf.Bytes())
	if decompSize > pbPool.limit {
		return nil, fmt.Errorf("decompressed request too long")
	}
	if err != nil {
		return nil, err
	}
	slice, err := pbPool.acquireSlice(decompSize)
	if err != nil {
		return nil, err
	}
	_, err = snappy.Decode(slice.Bytes(), buf.Bytes())
	return slice, err
}

func SetGlobalLimit(limit int) {
	requestPool.limit = limit / 2
	requestPool.rateLimiter = semaphore.NewWeighted(int64(limit / 2))
	pbPool.limit = limit / 2
	pbPool.rateLimiter = semaphore.NewWeighted(int64(limit / 2))
}
