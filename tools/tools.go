//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/ServiceWeaver/weaver/cmd/weaver"
	_ "github.com/amacneil/dbmate/v2"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "github.com/spf13/cobra-cli"
)
