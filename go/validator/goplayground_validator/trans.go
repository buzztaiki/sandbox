package goplayground_validator

import (
	"fmt"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/ja"
)

type translator struct {
	registerDefault func(*validator.Validate, ut.Translator) (err error)
	texts           map[string]string
}

func translators() map[string]translator {
	return map[string]translator{
		"en": {
			registerDefault: en_translations.RegisterDefaultTranslations,
			texts: map[string]string{
				"sex": "{0} should be male or female",
			},
		},
		"ja": {
			registerDefault: en_translations.RegisterDefaultTranslations,
			texts: map[string]string{
				"sex": "{0} には male または femail を指定してください",
			},
		},
	}
}

func registerTranslations(validate *validator.Validate, trans ut.Translator, texts map[string]string) error {
	for tag, text := range texts {
		err := validate.RegisterTranslation(
			tag, trans,
			func(ut ut.Translator) error { return ut.Add(tag, text, true) },
			validators()[tag].Translate,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetupTranslation(validate *validator.Validate) (*ut.UniversalTranslator, error) {
	locs := []locales.Translator{en.New(), ja.New()}
	uni := ut.New(locs[0], locs...)

	for locale, x := range translators() {
		trans, found := uni.GetTranslator(locale)
		if !found {
			return nil, fmt.Errorf("locale not found: %s", locale)
		}
		if err := x.registerDefault(validate, trans); err != nil {
			return nil, err
		}
		if err := registerTranslations(validate, trans, x.texts); err != nil {
			return nil, err
		}
	}

	return uni, nil
}
