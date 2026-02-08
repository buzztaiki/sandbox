package use_gotext

import (
	"github.com/leonelquinteros/gotext"
)

func HelloWorld(l *gotext.Locale) string {
	return l.Get("Hello, World!")
}

func Purchase(l *gotext.Locale, item, place string) string {
	return l.Get("Purchased %s at %s.", item, place)
}
