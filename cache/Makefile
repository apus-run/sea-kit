# 单元测试
.PHONY: ut
ut:
	@go test -race ./...

# 初始化环境
.PHONY: setup
setup:
	@sh ./script/setup.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	@sh ./script/fmt.sh

.PHONY: tidy
tidy:
	@go mod tidy -v

.PHONY: check
check:
	@$(MAKE) --no-print-directory fmt
	@$(MAKE) --no-print-directory tidy


# 生成 mock 文件
.PHONY: mock
mock:
	@mockgen -source=./cache.go -package=mocks -destination=./mocks/cache.mock.go
	@mockgen -package=mocks -destination=./mocks/cmd.mock.go github.com/redis/go-redis/v9 Cmdable

	@$(MAKE) tidy
	@echo "go mock finished"


# e2e 测试
.PHONY: e2e
e2e:
	sh ./script/integrate_test.sh