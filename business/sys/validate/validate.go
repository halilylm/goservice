package validate

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	tr_translations "github.com/go-playground/validator/v10/translations/tr"
	"reflect"
	"strings"
)

var validate *validator.Validate
var translator ut.Translator

func init() {
	validate = validator.New()

	translator, _ = ut.New(en.New(), tr.New()).GetTranslator("tr")

	tr_translations.RegisterDefaultTranslations(validate, translator)

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Check(s any) error {
	if err := validate.Struct(s); err != nil {
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(translator),
			}
			fields = append(fields, field)
		}
		return fields
	}
	return nil
}
