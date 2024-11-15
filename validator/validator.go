package validator

import (
	"context"
	"fmt"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhcn "github.com/go-playground/validator/v10/translations/zh"
)

var v *Validator

// Validator 可被用于Gin框架的验证器
// 具体支持的验证规则，可以参考：https://pkg.go.dev/github.com/go-playground/validator/v10
type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func (v Validator) Valid(obj any) error {
	return v.valid(obj)
}

func (v Validator) ValidCtx(ctx context.Context, obj any) error {
	return v.validCtx(ctx, obj)
}

// Engine 实现Gin验证器接口
func (v Validator) Engine() any {
	return v.validator
}

// New 生成一个验证器实例
// 在Gin中使用：binding.Validator = validator.New()
func New(opts ...Option) *Validator {
	validate := validator.New()
	validate.SetTagName("valid")

	zhTrans := zh.New()
	trans, _ := ut.New(zhTrans, zhTrans).GetTranslator("zh")

	if err := zhcn.RegisterDefaultTranslations(validate, trans); err != nil {
		fmt.Errorf("validator translation: %v", err)
	}

	for _, f := range opts {
		f(validate, trans)
	}

	return &Validator{
		validator:  validate,
		translator: trans,
	}
}

// InitValidator 默认初始化 validator.Validate
func InitValidator() {
	v = New()
}

// GetValidator .
func GetValidator() *Validator {
	return v
}

// V 是 GetValidator 简写
func V() *Validator {
	return GetValidator()
}

// ValidateStruct 验证结构体
func ValidateStruct(obj any) error {
	return v.Valid(obj)
}

// ValidateStructCtx 验证结构体，带Context
func ValidateStructCtx(ctx context.Context, obj any) error {
	return v.ValidCtx(ctx, obj)
}
