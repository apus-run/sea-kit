module github.com/apus-run/sea-kit/cache

go 1.21

require (
	github.com/apus-run/sea-kit/collection v0.0.0-00010101000000-000000000000
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/redis/go-redis/v9 v9.1.0
	github.com/smartystreets/goconvey v1.8.1
	github.com/stretchr/testify v1.8.4
	go.uber.org/mock v0.3.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/smarty/assertions v1.15.0 // indirect
	golang.org/x/exp v0.0.0-20240119083558-1b970713d09a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/apus-run/sea-kit/collection => ../collection
