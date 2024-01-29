module github.com/apus-run/sea-kit/cache

go 1.21

require (
	github.com/apus-run/sea-kit/utils v0.0.0-20231215063945-2c0bdf2b5759
	github.com/ecodeclub/ecache v0.0.0-20231031072032-8c6eedcb16de
	github.com/ecodeclub/ekit v0.0.8
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
	github.com/google/uuid v1.3.1 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/lithammer/shortuuid/v4 v4.0.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/smarty/assertions v1.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/apus-run/sea-kit/utils => ../utils
	github.com/apus-run/sea-kit/collection  => ../collection
)
