module github.com/apus-run/sea-kit/slogx

go 1.21

require github.com/natefinch/lumberjack v2.0.0+incompatible

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/apus-run/sea-kit/errorsx => ../errorsx
)

exclude github.com/mitchellh/osext v0.0.0-20151018003038-5e2d6d41470f