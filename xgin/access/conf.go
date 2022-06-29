package access

type BaseConf struct {
	RequestHeader  bool
	RequestBody    bool
	ResponseHeader bool
	ResponseBody   bool
}

type Conf struct {
	BaseConf
	SkipPaths    []string
	SpecificPath map[string]BaseConf
	Handler      func(entry *Entry)
}
