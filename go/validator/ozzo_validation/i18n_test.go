package ozzo_validation

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func ExampleNewI18nBundle() {
	bundle, err := NewI18nBundle()
	if err != nil {
		panic(err)
	}

	en := i18n.NewLocalizer(bundle, "en")
	fmt.Println(en.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    "HelloPerson",
		TemplateData: map[string]string{"Name": "John"},
	}))

	ja := i18n.NewLocalizer(bundle, "ja")
	fmt.Println(ja.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    "HelloPerson",
		TemplateData: map[string]string{"Name": "John"},
	}))

	// Output:
	// Hello John
	// こんにちは John
}
