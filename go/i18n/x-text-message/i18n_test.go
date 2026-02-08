package xtextmessage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/message"
)

func TestJa(t *testing.T) {
	lang := message.MatchLanguage("ja")
	require.Equal(t, "こんにちは世界！", HelloWorld(lang))
	require.Equal(t, "Amazonで本を購入しました。", Purchase(lang, "本", "Amazon"))
}

func TestEn(t *testing.T) {
	lang := message.MatchLanguage("en")
	require.Equal(t, "Hello, World!", HelloWorld(lang))
	require.Equal(t, "Purchased book at Amazon.", Purchase(lang, "book", "Amazon"))
}
