module github.com/buzztaiki/sandbox/go/i18n

go 1.25.7

require (
	github.com/leonelquinteros/gotext v1.7.2
	github.com/nicksnyder/go-i18n/v2 v2.6.1
	github.com/pelletier/go-toml/v2 v2.2.4
	github.com/stretchr/testify v1.11.1
	golang.org/x/text v0.33.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/tools v0.40.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

tool (
	github.com/leonelquinteros/gotext/cli/xgotext
	golang.org/x/text/cmd/gotext
)
