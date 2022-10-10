package bind

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"reflect"
)

var (
	ErrorCallback func(ctx *gin.Context, err error)
)

func bind(bindKey Key, val interface{}) gin.HandlerFunc {
	value := reflect.ValueOf(val)
	if value.Kind() == reflect.Ptr {
		panic(`Bind struct can not be a pointer. Example: Use: bind(Struct{}) instead of bind(&Struct{})`)
	}
	typ := value.Type()
	return func(ctx *gin.Context) {
		obj := reflect.New(typ).Interface()
		var err error
		switch bindKey {
		case KeyJSON:
			err = ctx.ShouldBindJSON(obj)
		case KeyXML:
			err = ctx.ShouldBindXML(obj)
		case KeyForm:
			err = ctx.ShouldBindWith(obj, binding.Form)
		case KeyQuery:
			err = ctx.ShouldBindQuery(obj)
		case KeyFormPost:
			err = ctx.ShouldBindWith(obj, binding.FormPost)
		case KeyFormMultipart:
			err = ctx.ShouldBindWith(obj, binding.FormMultipart)
		case KeyProtoBuf:
			err = ctx.ShouldBindWith(obj, binding.ProtoBuf)
		case KeyMsgPack:
			err = ctx.ShouldBindWith(obj, binding.MsgPack)
		case KeyYAML:
			err = ctx.ShouldBindYAML(obj)
		case KeyHeader:
			err = ctx.ShouldBindHeader(obj)
		case KeyTOML:
			err = ctx.ShouldBindTOML(obj)
		case KeyUri:
			err = ctx.ShouldBindUri(obj)
		}
		if v, ok := obj.(Validator); ok {
			err = v.Validate()
		}
		if err != nil {
			if ErrorCallback != nil {
				ErrorCallback(ctx, err)
				ctx.Abort()
			} else {
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			return
		}
		ctx.Set(string(bindKey), obj)
	}
}

func JSON(val interface{}) gin.HandlerFunc {
	return bind(KeyJSON, val)
}

func XML(val interface{}) gin.HandlerFunc {
	return bind(KeyXML, val)
}

func Form(val interface{}) gin.HandlerFunc {
	return bind(KeyForm, val)
}

func Query(val interface{}) gin.HandlerFunc {
	return bind(KeyQuery, val)
}

func FormPost(val interface{}) gin.HandlerFunc {
	return bind(KeyFormPost, val)
}

func FormMultipart(val interface{}) gin.HandlerFunc {
	return bind(KeyFormMultipart, val)
}

func ProtoBuf(val interface{}) gin.HandlerFunc {
	return bind(KeyProtoBuf, val)
}

func MsgPack(val interface{}) gin.HandlerFunc {
	return bind(KeyMsgPack, val)
}

func YAML(val interface{}) gin.HandlerFunc {
	return bind(KeyYAML, val)
}

func Header(val interface{}) gin.HandlerFunc {
	return bind(KeyHeader, val)
}

func TOML(val interface{}) gin.HandlerFunc {
	return bind(KeyTOML, val)
}

func Uri(val interface{}) gin.HandlerFunc {
	return bind(KeyUri, val)
}
