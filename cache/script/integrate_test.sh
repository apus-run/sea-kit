#!/usr/bin/env bash

set -e
docker-compose -f script/integration_test_docker_compose.yml down
docker-compose -f script/integration_test_docker_compose.yml up

go test -race -cover ./... -tags=e2e
docker-compose -f script/integration_test_docker_compose.yml down