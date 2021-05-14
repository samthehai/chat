// +build tools

package main

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/google/wire/cmd/wire"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/stringer"
)
