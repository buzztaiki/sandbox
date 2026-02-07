.PHONY: pot
pot:
	go run github.com/leonelquinteros/gotext/cli/xgotext@latest -in . -out ./po

.PHONY: msginit
msginit:
	msginit -i po/default.pot -l $(L) -o po/$(L)/default.po --no-translator

.PHONY: msgmerge
msgmerge:
	msgmerge -U po/$(L)/default.po po/default.pot
