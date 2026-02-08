package goi18n

import (
	"embed"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

//go:embed locale.*.toml
var LocaleFS embed.FS

func loadBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(LocaleFS, "locale.en.toml")
	bundle.LoadMessageFileFS(LocaleFS, "locale.ja.toml")
	return bundle
}

func TestJa(t *testing.T) {
	bundle := loadBundle()
	localizer := i18n.NewLocalizer(bundle, "ja")

	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: "helloWorld"})
	require.NoError(t, err)
	require.Equal(t, "こんにちは世界！", msg)

	msg, err = localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "purchase",
		TemplateData: map[string]any{"location": "Amazon", "product": "本"}})
	require.NoError(t, err)
	require.Equal(t, "Amazonで本を購入しました。", msg)
}

func TestEn(t *testing.T) {
	bundle := loadBundle()
	localizer := i18n.NewLocalizer(bundle, "en")

	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: "helloWorld"})
	require.NoError(t, err)
	require.Equal(t, "Hello, world!", msg)

	msg, err = localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "purchase",
		TemplateData: map[string]any{"location": "Amazon", "product": "book"}})
	require.NoError(t, err)
	require.Equal(t, "Purchased book at Amazon.", msg)
}
