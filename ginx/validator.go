package ginx

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	tzh "github.com/go-playground/validator/v10/translations/zh"

	"github.com/apus-run/sea-kit/ginx/validators"
)

func init() {
	binding.Validator = &defaultValidator{}
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj any) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		v.lazyInit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

func newValidator() *validator.Validate {
	// 注册translator
	zhTranslator := zh.New()
	uni := ut.New(zhTranslator, zhTranslator)
	trans, _ = uni.GetTranslator("zh")
	v := validator.New()
	v.RegisterValidation("notBlank", validators.NotBlank)
	v.RegisterValidation("email", validators.ValidEmail)
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("yaml"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	if err := tzh.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal("Gin fail to registered Translation")
	}
	return v
}

func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = newValidator()
		v.validate.SetTagName("binding")
	})
}

var trans ut.Translator

func validate(errs error) error {
	if validationErrors, ok := errs.(validator.ValidationErrors); ok {
		var errList []string
		for _, e := range validationErrors {
			errList = append(errList, e.Translate(trans))
		}
		return errors.New(strings.Join(errList, "|"))
	}
	return errs
}
