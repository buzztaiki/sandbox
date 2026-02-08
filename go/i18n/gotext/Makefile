.PHONY: pot
pot:
	go tool xgotext -in . -out ./po

.PHONY: msginit
msginit:
	msginit -i po/default.pot -l $(L) -o po/$(L)/default.po --no-translator

.PHONY: msgmerge
msgmerge:
	msgmerge -U po/$(L)/default.po po/default.pot
