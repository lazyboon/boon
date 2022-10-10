package xredis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"strconv"
	"time"
)

var (
	// luaLockReleaseScript
	// KEYS[1] - lock name
	// ARGV[1] - token
	// return 1 if the lock was released, otherwise 0
	luaLockReleaseScript = redis.NewScript(`
		if redis.call('get', KEYS[1]) == ARGV[1] then 
			return redis.call('del', KEYS[1]) 
		else 
			return 0 
		end
	`)

	// luaLockRefreshScript
	// # KEYS[1] - lock name
	// # ARGV[1] - token
	// # ARGV[2] - milliseconds
	// # return 1 if the locks time was reacquired, otherwise 0
	luaLockRefreshScript = redis.NewScript(`
		if redis.call('get', KEYS[1]) == ARGV[1] then 
			return redis.call('pexpire', KEYS[1], ARGV[2]) 
		else 
			return 0 
		end
	`)
)

var (
	// ErrAcquireLock is returned when a lock can't acquire.
	ErrAcquireLock = errors.New("redis lock: failed to acquire the lock")

	// ErrLockInactive is returned when trying to release an inactive lock.
	ErrLockInactive = errors.New("redis lock: lock inactive")
)

type Lock struct {
	client     *redis.Client
	key        string
	val        string
	expiration time.Duration
	token      string
}

func NewLock(ctx context.Context, client *redis.Client, key string, expiration time.Duration, options ...LockOption) (*Lock, error) {
	l := &Lock{
		client: client,
	}
	// new mutex lock options
	opts := newLockOptions(options...)

	// assignment
	l.key = key
	l.expiration = expiration

	// generate a unique key to prevent others from releasing the lock
	l.token = uuid.New().String()
	l.val = opts.value

	// timeout
	timeout := expiration
	if opts.blockingTimeout != nil {
		timeout = *opts.blockingTimeout
	}

	// make sure can exit
	deadline := time.Now().Add(timeout)
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	value := fmt.Sprintf("%s%s", l.token, l.val)
	var doCount uint
	for {
		ok, err := l.client.SetNX(ctx, key, value, expiration).Result()
		if err != nil {
			return nil, err
		}
		if ok {
			return l, nil
		}
		doCount++
		stop, duration := opts.backoff.Next(doCount)
		if stop {
			return nil, ErrAcquireLock
		}
		if time.Now().Add(duration).After(deadline) {
			return nil, ErrAcquireLock
		}
		timer.Reset(duration)
		select {
		case <-ctx.Done():
			return nil, ErrAcquireLock
		case <-timer.C:
		}
	}
}

func (l *Lock) Key() string {
	return l.key
}

func (l *Lock) Value() string {
	return l.val
}

func (l *Lock) Token() string {
	return l.token
}

func (l *Lock) Release(ctx context.Context) error {
	res, err := luaLockReleaseScript.Run(ctx, l.client, []string{l.key}, fmt.Sprintf("%s%s", l.token, l.val)).Result()
	if err == redis.Nil {
		return ErrLockInactive
	}
	if err != nil {
		return err
	}
	if i, ok := res.(int64); !ok || i != 1 {
		return ErrLockInactive
	}
	return nil
}

func (l *Lock) Refresh(ctx context.Context, expiration time.Duration) error {
	milliseconds := strconv.FormatInt(int64(expiration/time.Millisecond), 10)
	value := fmt.Sprintf("%s%s", l.token, l.val)
	status, err := luaLockRefreshScript.Run(ctx, l.client, []string{l.key}, value, milliseconds).Result()
	if err != nil {
		return err
	}
	if status != 1 {
		return ErrAcquireLock
	}
	return nil
}
