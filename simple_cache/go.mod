module github.com/apus-run/sea-kit/simple_cache

go 1.20

replace (
	github.com/apus-run/sea-kit/list => ../list
	github.com/apus-run/sea-kit/set => ../set
	github.com/apus-run/sea-kit/utils => ../utils
)

require (
	github.com/apus-run/sea-kit/list v0.0.0-00010101000000-000000000000
	github.com/apus-run/sea-kit/set v0.0.0-00010101000000-000000000000
	github.com/apus-run/sea-kit/utils v0.0.0-00010101000000-000000000000
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/redis/go-redis/v9 v9.3.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/mock v0.3.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
