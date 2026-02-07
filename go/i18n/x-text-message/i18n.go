package xtextmessage

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func HelloWorld(lang language.Tag) string {
	return message.NewPrinter(lang).Sprintf("Hello, World!")
}

func Purchase(lang language.Tag, item, place string) string {
	return message.NewPrinter(lang).Sprintf("Purchased %s at %s.", item, place)
}
