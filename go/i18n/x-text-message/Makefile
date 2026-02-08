update:
	go tool gotext update -lang=en-US,ja-JP -out=locales/catalog_gen.go
	find locales -name 'out.gotext.json' | xargs -I% bash -c 'cd $$(dirname $$0) && mv -f out.gotext.json messages.gotext.json' %
