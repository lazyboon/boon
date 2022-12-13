package access

type BaseOption struct {
	RequestHeader  *bool
	RequestBody    *bool
	ResponseHeader *bool
	ResponseBody   *bool
}

func NewBaseOption() *BaseOption {
	return &BaseOption{}
}

func (b *BaseOption) SetRequestHeader(v bool) *BaseOption {
	b.RequestHeader = &v
	return b
}

func (b *BaseOption) SetRequestBody(v bool) *BaseOption {
	b.RequestBody = &v
	return b
}

func (b *BaseOption) SetResponseHeader(v bool) *BaseOption {
	b.ResponseHeader = &v
	return b
}

func (b *BaseOption) SetResponseBody(v bool) *BaseOption {
	b.ResponseBody = &v
	return b
}

type Option struct {
	*BaseOption
	SkipPaths    []string
	SpecificPath map[string]*BaseOption
	Handler      func(entity *Entity)
}

func NewOption() *Option {
	return &Option{
		BaseOption: &BaseOption{
			RequestHeader:  nil,
			RequestBody:    nil,
			ResponseHeader: nil,
			ResponseBody:   nil,
		},
		SkipPaths:    make([]string, 0),
		SpecificPath: map[string]*BaseOption{},
		Handler:      nil,
	}
}

func (o *Option) SetBaseOption(base *BaseOption) *Option {
	o.BaseOption = base
	return o
}

func (o *Option) SetSkipPaths(v []string) *Option {
	o.SkipPaths = v
	return o
}

func (o *Option) SetSpecificPath(v map[string]*BaseOption) *Option {
	o.SpecificPath = v
	return o
}

func (o *Option) SetHandler(v func(entity *Entity)) *Option {
	o.Handler = v
	return o
}

func mergeOptions(options ...*Option) *Option {
	ans := NewOption()
	for _, item := range options {
		if item.BaseOption != nil {
			if item.BaseOption.RequestHeader != nil {
				ans.BaseOption.RequestHeader = item.BaseOption.RequestHeader
			}
			if item.BaseOption.RequestBody != nil {
				ans.BaseOption.RequestBody = item.BaseOption.RequestBody
			}
			if item.BaseOption.ResponseHeader != nil {
				ans.BaseOption.ResponseHeader = item.BaseOption.ResponseHeader
			}
			if item.BaseOption.ResponseBody != nil {
				ans.BaseOption.ResponseBody = item.BaseOption.ResponseBody
			}
		}
		if item.SkipPaths != nil {
			ans.SkipPaths = append(ans.SkipPaths, item.SkipPaths...)
		}
		for key, val := range item.SpecificPath {
			ans.SpecificPath[key] = val
		}
		if item.Handler != nil {
			ans.Handler = item.Handler
		}
	}
	return ans
}
