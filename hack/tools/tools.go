// +build tools

package tools

// Mock imports to enforce their installation by `go mod`.
// Go modules will be forced to download and install them.
// See: https://git.io/Jewio and https://git.io/Je7Xh.

import (
	_ "github.com/go-bindata/go-bindata"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
)
