// +build tools

package tools

import (
    _ "github.com/cweill/gotests/gotests"
    _ "github.com/fatih/gomodifytags"
    _ "github.com/josharian/impl"
    _ "github.com/haya14busa/goplay"
    _ "github.com/go-delve/delve/cmd/dlv"
    _ "honnef.co/go/tools/cmd/staticcheck"
    _ "golang.org/x/tools/gopls"
)

