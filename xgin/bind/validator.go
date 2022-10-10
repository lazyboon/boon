package bind

type Validator interface {
	Validate() error
}
