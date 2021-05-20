// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/rubenv/sql-migrate/sql-migrate"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/stringer"
)
