//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/amacneil/dbmate/v2"
	_ "github.com/spf13/cobra-cli"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
