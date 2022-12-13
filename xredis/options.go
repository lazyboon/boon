package xredis

import (
	"github.com/lazyboon/boon/retry"
	"time"
)

type LockOption struct {
	Value           *string
	BlockingTimeout *time.Duration
	Backoff         retry.IBackoff
}

func NewLockOption() *LockOption {
	return &LockOption{
		Backoff: &retry.ZeroBackoff{},
	}
}

func (l *LockOption) SetValue(v string) *LockOption {
	l.Value = &v
	return l
}

func (l *LockOption) SetBlockingTimeout(v time.Duration) *LockOption {
	l.BlockingTimeout = &v
	return l
}

func (l *LockOption) SetBackoff(v retry.IBackoff) *LockOption {
	l.Backoff = v
	return l
}

func mergeLockOptions(options ...*LockOption) *LockOption {
	ans := NewLockOption()
	for _, item := range options {
		if item.Value != nil {
			ans.Value = item.Value
		}
		if item.BlockingTimeout != nil {
			ans.BlockingTimeout = item.BlockingTimeout
		}
		if item.Backoff != nil {
			ans.Backoff = item.Backoff
		}
	}
	return ans
}

//----------------------------------------------------------------------------------------------------------------------

type DelayOption struct {
	Namespace     *string
	ListenTopics  []string
	ErrorCallback func(err error)
}

func NewDelayOption() *DelayOption {
	namespace := "com.lazyboon"
	return &DelayOption{
		Namespace: &namespace,
	}
}

func (d *DelayOption) SetNamespace(v string) *DelayOption {
	d.Namespace = &v
	return d
}

func (d *DelayOption) SetListenTopics(v []string) *DelayOption {
	d.ListenTopics = v
	return d
}

func (d *DelayOption) SetErrorCallback(v func(err error)) *DelayOption {
	d.ErrorCallback = v
	return d
}

func mergeDelayOption(options ...*DelayOption) *DelayOption {
	ans := NewDelayOption()
	for _, item := range options {
		if item.Namespace != nil {
			ans.Namespace = item.Namespace
		}
		if item.ListenTopics != nil {
			ans.ListenTopics = item.ListenTopics
		}
		if item.ErrorCallback != nil {
			ans.ErrorCallback = item.ErrorCallback
		}
	}
	return ans
}
