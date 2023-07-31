package xredis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
)

var (
	ErrCaptchaLimit = errors.New("captcha limit exceeded")
)

var (
	luaCaptchaSetScript = redis.NewScript(`
		local key = KEYS[1]
		local code = ARGV[1]
		local expired = ARGV[2]
		for k, v in pairs(ARGV) do
		  if k > 2 and k % 2 == 1 then
			local time_key = key .. ':' .. v
			redis.call('SET', time_key, 0, 'EX', tonumber(v), 'NX')
			local count = redis.call('INCR', time_key)
			local limit = tonumber(ARGV[k+1])
			if tonumber(count) > limit then
			  return {-1, tonumber(v), limit}
			end
		  end
		end
		redis.call('SET', key .. ':key', code, 'EX', tonumber(expired))
		return {1, 0, 0}
	`)
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ICaptcha interface {
	Get(key string) (string, error)
	Set(key string, value string) (*CaptchaRate, error)
	Delete(key string) error
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Captcha struct {
	client *Client
	option *CaptchaOption
}

func NewCaptcha(client *Client, options ...*CaptchaOption) (*Captcha, error) {
	opts := mergeCaptchaOption(options...)
	if len(opts.Rates) == 0 {
		return nil, errors.New("captcha rates must provide")
	}
	return &Captcha{
		client: client,
		option: opts,
	}, nil
}

func (c *Captcha) Get(key string) (string, error) {
	return c.client.Get(context.TODO(), fmt.Sprintf("%s:key", c.key(key))).Result()
}

func (c *Captcha) Set(key string, value string) (*CaptchaRate, error) {
	argv := []interface{}{value, *c.option.Expire}
	for _, item := range c.option.Rates {
		argv = append(argv, item.Seconds, item.Count)
	}
	rsp, err := luaCaptchaSetScript.Run(context.TODO(), c.client, []string{c.key(key)}, argv...).Uint64Slice()
	if err != nil {
		return nil, err
	}
	if rsp[0] == 1 {
		return nil, nil
	}
	return &CaptchaRate{
		Seconds: uint(rsp[1]),
		Count:   uint(rsp[2]),
	}, ErrCaptchaLimit
}

func (c *Captcha) Delete(key string) error {
	_, err := c.client.Del(context.TODO(), fmt.Sprintf("%s:key", c.key(key))).Result()
	return err
}

func (c *Captcha) key(key string) string {
	return fmt.Sprintf("%s:%s", *c.option.Namespace, key)
}
