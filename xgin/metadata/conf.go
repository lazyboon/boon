package metadata

type Option struct {
	RequestID     *bool
	ReceiveTime   *bool
	ResponseTime  *bool
	ServerName    *string
	ServerVersion *string
}

func NewOption() *Option {
	return &Option{}
}

func (o *Option) SetRequestID(b bool) *Option {
	o.RequestID = &b
	return o
}

func (o *Option) SetReceiveTime(b bool) *Option {
	o.ReceiveTime = &b
	return o
}

func (o *Option) SetResponseTime(b bool) *Option {
	o.ResponseTime = &b
	return o
}

func (o *Option) SetServerName(name string) *Option {
	o.ServerName = &name
	return o
}

func (o *Option) SetServerVersion(version string) *Option {
	o.ServerVersion = &version
	return o
}

func mergeOptions(options ...*Option) *Option {
	ans := NewOption()
	for _, item := range options {
		if item.RequestID != nil {
			ans.RequestID = item.RequestID
		}
		if item.ReceiveTime != nil {
			ans.ReceiveTime = item.ReceiveTime
		}
		if item.ResponseTime != nil {
			ans.ResponseTime = item.ResponseTime
		}
		if item.ServerName != nil {
			ans.ServerName = item.ServerName
		}
		if item.ServerVersion != nil {
			ans.ServerVersion = item.ServerVersion
		}
	}
	return ans
}
