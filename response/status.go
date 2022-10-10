package response

import (
	"net/http"
)

// http status code 1xx
var (
	Continue           = NewWithStatusCode(http.StatusContinue)
	SwitchingProtocols = NewWithStatusCode(http.StatusSwitchingProtocols)
	Processing         = NewWithStatusCode(http.StatusProcessing)
	EarlyHints         = NewWithStatusCode(http.StatusEarlyHints)
)

// http status code 2xx
var (
	OK                   = NewWithStatusCode(http.StatusOK)
	Created              = NewWithStatusCode(http.StatusCreated)
	Accepted             = NewWithStatusCode(http.StatusAccepted)
	NonAuthoritativeInfo = NewWithStatusCode(http.StatusNonAuthoritativeInfo)
	NoContent            = NewWithStatusCode(http.StatusNoContent)
	ResetContent         = NewWithStatusCode(http.StatusResetContent)
	PartialContent       = NewWithStatusCode(http.StatusPartialContent)
	MultiStatus          = NewWithStatusCode(http.StatusMultiStatus)
	AlreadyReported      = NewWithStatusCode(http.StatusAlreadyReported)
	IMUsed               = NewWithStatusCode(http.StatusIMUsed)
)

// http status code 3xx
var (
	MultipleChoices   = NewWithStatusCode(http.StatusMultipleChoices)
	MovedPermanently  = NewWithStatusCode(http.StatusMovedPermanently)
	Found             = NewWithStatusCode(http.StatusFound)
	SeeOther          = NewWithStatusCode(http.StatusSeeOther)
	NotModified       = NewWithStatusCode(http.StatusNotModified)
	UseProxy          = NewWithStatusCode(http.StatusUseProxy)
	TemporaryRedirect = NewWithStatusCode(http.StatusTemporaryRedirect)
	PermanentRedirect = NewWithStatusCode(http.StatusPermanentRedirect)
)

// http status code 4xx
var (
	BadRequest                   = NewWithStatusCode(http.StatusBadRequest)
	Unauthorized                 = NewWithStatusCode(http.StatusUnauthorized)
	PaymentRequired              = NewWithStatusCode(http.StatusPaymentRequired)
	Forbidden                    = NewWithStatusCode(http.StatusForbidden)
	NotFound                     = NewWithStatusCode(http.StatusNotFound)
	MethodNotAllowed             = NewWithStatusCode(http.StatusMethodNotAllowed)
	NotAcceptable                = NewWithStatusCode(http.StatusNotAcceptable)
	ProxyAuthRequired            = NewWithStatusCode(http.StatusProxyAuthRequired)
	RequestTimeout               = NewWithStatusCode(http.StatusRequestTimeout)
	Conflict                     = NewWithStatusCode(http.StatusConflict)
	Gone                         = NewWithStatusCode(http.StatusGone)
	LengthRequired               = NewWithStatusCode(http.StatusLengthRequired)
	PreconditionFailed           = NewWithStatusCode(http.StatusPreconditionFailed)
	RequestEntityTooLarge        = NewWithStatusCode(http.StatusRequestEntityTooLarge)
	RequestURITooLong            = NewWithStatusCode(http.StatusRequestURITooLong)
	UnsupportedMediaType         = NewWithStatusCode(http.StatusUnsupportedMediaType)
	RequestedRangeNotSatisfiable = NewWithStatusCode(http.StatusRequestedRangeNotSatisfiable)
	ExpectationFailed            = NewWithStatusCode(http.StatusExpectationFailed)
	Teapot                       = NewWithStatusCode(http.StatusTeapot)
	MisdirectedRequest           = NewWithStatusCode(http.StatusMisdirectedRequest)
	UnprocessableEntity          = NewWithStatusCode(http.StatusUnprocessableEntity)
	Locked                       = NewWithStatusCode(http.StatusLocked)
	FailedDependency             = NewWithStatusCode(http.StatusFailedDependency)
	TooEarly                     = NewWithStatusCode(http.StatusTooEarly)
	UpgradeRequired              = NewWithStatusCode(http.StatusUpgradeRequired)
	PreconditionRequired         = NewWithStatusCode(http.StatusPreconditionRequired)
	TooManyRequests              = NewWithStatusCode(http.StatusTooManyRequests)
	RequestHeaderFieldsTooLarge  = NewWithStatusCode(http.StatusRequestHeaderFieldsTooLarge)
	UnavailableForLegalReasons   = NewWithStatusCode(http.StatusUnavailableForLegalReasons)
)

// http status code 5xx
var (
	InternalServerError           = NewWithStatusCode(http.StatusInternalServerError)
	NotImplemented                = NewWithStatusCode(http.StatusNotImplemented)
	BadGateway                    = NewWithStatusCode(http.StatusBadGateway)
	ServiceUnavailable            = NewWithStatusCode(http.StatusServiceUnavailable)
	GatewayTimeout                = NewWithStatusCode(http.StatusGatewayTimeout)
	HTTPVersionNotSupported       = NewWithStatusCode(http.StatusHTTPVersionNotSupported)
	VariantAlsoNegotiates         = NewWithStatusCode(http.StatusVariantAlsoNegotiates)
	InsufficientStorage           = NewWithStatusCode(http.StatusInsufficientStorage)
	LoopDetected                  = NewWithStatusCode(http.StatusLoopDetected)
	NotExtended                   = NewWithStatusCode(http.StatusNotExtended)
	NetworkAuthenticationRequired = NewWithStatusCode(http.StatusNetworkAuthenticationRequired)
)
