package ozzo_validation

import (
	"embed"
	"io/fs"
	"log/slog"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goccy/go-yaml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed i18n/*.yaml
var i18nFS embed.FS

func NewI18nBundle() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	matches, err := fs.Glob(i18nFS, "i18n/*.yaml")
	if err != nil {
		return nil, err
	}
	for _, path := range matches {
		_, err := bundle.LoadMessageFileFS(i18nFS, path)
		if err != nil {
			return nil, err
		}
	}
	return bundle, nil
}

func TranslateErrors(errs validation.Errors, localizer *i18n.Localizer) {
	for k, e := range errs {
		switch err := e.(type) {
		case validation.Errors:
			TranslateErrors(err, localizer)
		case validation.Error:
			message, lerr := localizer.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    err.Code(),
					Other: err.Message(),
				},
				TemplateData: err.Params(),
			})
			if lerr != nil {
				slog.Warn("failed to localize error message", "error", lerr, "code", err.Code(), "message", err.Message(), "params", err.Params())
				continue
			}
			errs[k] = err.SetMessage(message)
		}
	}
}
