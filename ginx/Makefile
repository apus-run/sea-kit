.PHONY: mock
mock:
	@mockgen -package=redismocks -destination ./ginx/middleware/activelimit/mocks/cmd.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy