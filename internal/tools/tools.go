//go:build tools
// +build tools

package tools

import (
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/goreleaser/nfpm"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/stringer"
	_ "google.golang.org/protobuf/protoc-gen-go"
)
