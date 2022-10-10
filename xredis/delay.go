package xredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"math/rand"
	"time"
)

var (
	ErrDelayQueueListenTopic     = errors.New("listen topic error")
	ErrDelayQueueListenReady     = errors.New("listen ready error")
	ErrDelayQueueGetJobFromPool  = errors.New("get job from pool error")
	ErrDelayQueueGetJobUnmarshal = errors.New("get job unmarshal error")
)

var (
	// luaDelayQueueAddJobScript
	// KEYS[1] - job key
	// ARGV[1] - job content json string
	// KEYS[2] - topic key
	// ARGV[2] - execution time, second timestamp
	// ARGV[3] - job id
	luaDelayQueueAddJobScript = redis.NewScript(`
		for k, v in pairs(ARGV) do
			local m = (k-1) % 3
			local n = math.floor((k - 1) / 3)
			if m == 0 then
				redis.call('SET', KEYS[n*2+1], v)
			elseif m == 1 then
				redis.call('ZADD', KEYS[n*2+2], v, ARGV[k+1])
			end
		end
		return 1
	`)

	// luaDelayQueueMoveToReadyScript
	// KEYS[1] - topic key
	// KEYS[2] - ready queue key
	// KEYS[3] - error jobs key
	// ARGV[1] - max score, second timestamp
	luaDelayQueueMoveToReadyScript = redis.NewScript(`
		local data = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1])
		local size = 0
		for k, v in pairs(data) do
			local job = redis.call('GET', v)
			if job ~= nil and job ~= false then
				local obj = cjson.decode(job)
				if obj.done > obj.retry then
					redis.call('SADD', KEYS[3], v)
					redis.call('ZREM', KEYS[1], v)
				else
					obj.done = obj.done + 1
					redis.call('ZADD', KEYS[1], obj.delay+(obj.ttr * obj.done), v)
					redis.call('SET', v, cjson.encode(obj))
					redis.call('LPUSH', KEYS[2], v)
				end
				size = size + 1
			else
				redis.call('ZREM', KEYS[1], v)
			end
		end
		return size
	`)

	// luaDelayQueueDeleteJobScript
	// KEYS[1] - topic key
	// KEYS[2] - job key
	luaDelayQueueDeleteJobScript = redis.NewScript(`
		redis.call('ZREM', KEYS[1], KEYS[2])
		redis.call('DEL', KEYS[2])
		return 1
	`)

	// luaDelayQueueDeleteErrorJobScript
	// KEYS[1] - error jobs key
	// KEYS[2] - job key
	luaDelayQueueDeleteErrorJobScript = redis.NewScript(`
		redis.call('SREM', KEYS[1], KEYS[2])
		redis.call('DEL', KEYS[2])
		return 1
	`)
)

type Job struct {
	Topic string `json:"topic"`
	ID    string `json:"id"`
	Delay int64  `json:"delay"`
	Body  string `json:"body"`
	Retry int64  `json:"retry"`
	TTR   int64  `json:"ttr"`
}

type jobWrapper struct {
	Job
	Done int64 `json:"done"`
}

type Delay struct {
	client    *Client
	readyChan chan *Job
	*delayOptions
}

func NewDelay(client *Client, options ...DelayOption) *Delay {
	rand.Seed(time.Now().UnixNano())
	opts := newDelayOptions(options...)
	d := &Delay{
		client:       client,
		delayOptions: opts,
		readyChan:    make(chan *Job),
	}
	d.listenZSet()
	d.listenReadyList()
	return d
}

func (d *Delay) Publish(jobs ...*Job) error {
	size := len(jobs)
	keys := make([]string, 0, size)
	argv := make([]interface{}, 0, size)
	for _, job := range jobs {
		if job.Topic == "" {
			return errors.New("topic cannot be empty")
		}
		if job.ID == "" {
			return errors.New("id cannot be empty")
		}
		if job.TTR <= 0 {
			return errors.New("ttr must be greater than 0")
		}
		raw, err := json.Marshal(jobWrapper{
			Job:  *job,
			Done: 0,
		})
		if err != nil {
			return err
		}
		poolJobKey := d.poolJobStringKey(job.Topic, job.ID)
		keys = append(keys, poolJobKey, d.topicZSetKey(job.Topic))
		argv = append(argv, string(raw), job.Delay, poolJobKey)
	}
	_, err := luaDelayQueueAddJobScript.Run(context.TODO(), d.client, keys, argv...).Result()
	if err != nil {
		return err
	}
	return nil
}

func (d *Delay) Commit(topic string, id string) error {
	_, err := luaDelayQueueDeleteJobScript.Run(context.TODO(), d.client, []string{d.topicZSetKey(topic), d.poolJobStringKey(topic, id)}).Result()
	return err
}

func (d *Delay) Ready() <-chan *Job {
	return d.readyChan
}

func (d *Delay) RandomGetErrorJob(topic string) (*Job, error) {
	data, err := d.client.SRandMemberN(context.TODO(), d.topicErrorSetKey(topic), 1).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, redis.Nil
	}
	raw, err := d.client.Get(context.TODO(), data[0]).Bytes()
	job := &jobWrapper{}
	err = json.Unmarshal(raw, job)
	if err != nil {
		return nil, err
	}
	return &job.Job, nil
}

func (d *Delay) RemoveErrorJob(topic string, id string) error {
	_, err := luaDelayQueueDeleteErrorJobScript.Run(context.TODO(), d.client, []string{d.topicErrorSetKey(topic), d.poolJobStringKey(topic, id)}).Result()
	return err
}

func (d *Delay) poolJobStringKey(topic string, id string) string {
	return fmt.Sprintf("%s:pool_%s:%s", d.namespace, topic, id)
}

func (d *Delay) topicReadyListKey(topic string) string {
	return fmt.Sprintf("%s:ready_%s", d.namespace, topic)
}

func (d *Delay) topicZSetKey(topic string) string {
	return fmt.Sprintf("%s:%s", d.namespace, topic)
}

func (d *Delay) topicErrorSetKey(topic string) string {
	return fmt.Sprintf("%s:error_%s", d.namespace, topic)
}

func (d *Delay) listenZSet() {
	for _, topic := range d.listenTopics {
		go d.listenZSetTopic(topic)
	}
}

func (d *Delay) listenZSetTopic(topic string) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	working := false
	for range ticker.C {
		if working {
			continue
		}
		working = true
		keys := []string{d.topicZSetKey(topic), d.topicReadyListKey(topic), d.topicErrorSetKey(topic)}
		_, err := luaDelayQueueMoveToReadyScript.Run(context.TODO(), d.client, keys, fmt.Sprintf("%d", time.Now().Unix())).Result()
		if err != nil && d.errorCallback != nil {
			d.errorCallback(fmt.Errorf("%w, %s", ErrDelayQueueListenTopic, err.Error()))
		}
		working = false
	}
}

func (d *Delay) listenReadyList() {
	go func() {
		topics := make([]string, len(d.listenTopics))
		for idx, topic := range d.listenTopics {
			topics[idx] = d.topicReadyListKey(topic)
		}
		for {
			result, err := d.client.BRPop(context.TODO(), time.Duration(10)*time.Second, topics...).Result()
			if err != nil {
				if !errors.Is(redis.Nil, err) && d.errorCallback != nil {
					d.errorCallback(fmt.Errorf("%w, %s", ErrDelayQueueListenReady, err.Error()))
				}
				continue
			}
			job, err := d.getJob(result[1])
			if err != nil && d.errorCallback != nil {
				d.errorCallback(err)
				continue
			}
			d.readyChan <- job
		}
	}()
}

func (d *Delay) getJob(key string) (*Job, error) {
	raw, err := d.client.Get(context.TODO(), key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrDelayQueueGetJobFromPool, err.Error())
	}
	job := &jobWrapper{}
	err = json.Unmarshal(raw, job)
	if err != nil {
		return nil, fmt.Errorf("%s, %s", ErrDelayQueueGetJobUnmarshal, err.Error())
	}
	return &job.Job, nil
}
