goimports -l -w $(find . -type f -name '*.go' -not -path "./.idea/*" -not -path "./vendor/*" -not -path "./.vscode/*")
