package xredis

import (
	"fmt"
	"math/rand"
	"time"
)

type Job struct {
	ID            string   `json:"id"`
	ExecTimestamp int64    `json:"exec_timestamp"`
	Payload       string   `json:"payload"`
	Backoff       IBackoff `json:"backoff"`
}

type delayQueue struct {
	bucketPrefix     string
	bucketCount      uint64
	scheduleInterval time.Duration
}

func NewDelayQueue(options ...DelayQueueOption) *delayQueue {
	rand.Seed(time.Now().UnixNano())
	opts := newDelayQueueOptions(options...)
	return &delayQueue{
		bucketPrefix:     opts.bucketPrefix,
		bucketCount:      opts.bucketCount,
		scheduleInterval: opts.scheduleInterval,
	}
}

// randomGetBucket randomly obtain a bucket according to the number of buckets in the configuration
func (d *delayQueue) randomGetBucket() string {
	return fmt.Sprintf("%s:%d", d.bucketPrefix, rand.Uint64()%d.bucketCount)
}

func (d *delayQueue) Upsert(job *Job) error {
	return nil
}

func (d *delayQueue) Delete(jobID string) error {
	return nil
}
