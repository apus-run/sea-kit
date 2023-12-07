package validator

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/text/language"
)

const I18nKey = "SpectatorNan/validate/i18n"

func GetLanguageTag(r *http.Request) language.Tag {
	accept := r.Header.Get("Accept-Language")
	langTags, _, err := language.ParseAcceptLanguage(accept)
	if err != nil {
		langTags = []language.Tag{language.English}
	}
	tags := []language.Tag{
		language.English,
		language.Spanish,
		language.Chinese,
	}
	var matcher = language.NewMatcher(tags)
	_, i, _ := matcher.Match(langTags...)
	//_, i := language.MatchStrings(matcher, langTag.String())
	tag := tags[i]
	return tag
}

func (v Validator) valid(obj any) error {
	if reflect.Indirect(reflect.ValueOf(obj)).Kind() != reflect.Struct {
		return nil
	}

	e := v.validator.Struct(obj)
	if e != nil {
		err, ok := e.(validator.ValidationErrors)
		if !ok {
			return e
		}
		return removeStructName(err.Translate(v.translator))
	}
	return nil
}

func (v Validator) validCtx(ctx context.Context, obj any) error {
	if reflect.Indirect(reflect.ValueOf(obj)).Kind() != reflect.Struct {
		return nil
	}

	e := v.validator.StructCtx(ctx, obj)
	if e != nil {
		err, ok := e.(validator.ValidationErrors)
		if !ok {
			return e
		}
		return removeStructName(err.Translate(v.translator))
	}
	return nil
}

func removeStructName(fields map[string]string) error {
	errs := make([]string, 0, len(fields))
	for _, err := range fields {
		errs = append(errs, err)
	}
	return errors.New(strings.Join(errs, ";"))
}
