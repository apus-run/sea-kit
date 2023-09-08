module github.com/apus-run/sea-kit/config

go 1.19

replace github.com/imdario/mergo => dario.cat/mergo v1.0.0

require (
	github.com/apus-run/sea-kit/encoding v0.0.0-20230905132113-81a1479e06c3
	github.com/apus-run/sea-kit/log v0.0.0-20230905132113-81a1479e06c3
	github.com/imdario/mergo v1.0.0
	google.golang.org/protobuf v1.31.0
)

require gopkg.in/yaml.v3 v3.0.1 // indirect
