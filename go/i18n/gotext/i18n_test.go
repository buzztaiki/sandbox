package use_gotext

import (
	"embed"
	"testing"

	"github.com/leonelquinteros/gotext"
	"github.com/stretchr/testify/require"
)

//go:embed po/*/*.po
var LocaleFS embed.FS

func TestJa(t *testing.T) {
	l := gotext.NewLocaleFSWithPath("ja_JP", LocaleFS, "po")
	l.AddDomain("default")
	require.Equal(t, "こんにちは世界！", HelloWorld(l))
	require.Equal(t, "Amazonで本を購入しました。", Purchase(l, "本", "Amazon"))
}

func TestEn(t *testing.T) {
	l := gotext.NewLocaleFSWithPath("ja_JP", LocaleFS, "po")
	l.AddDomain("default")
	require.Equal(t, "Hello, World!", HelloWorld(l))
	require.Equal(t, "Purchased book at Amazon.", Purchase(l, "book", "Amazon"))
}
