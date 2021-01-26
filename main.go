//go:generate trs proto/app.proto --lib=./proto --doc=./doc
//go:generate trs proto/share.proto --lib=./proto --doc=./doc
package main

import (
	"glab.tagtic.cn/ad_gains/kitty/cmd"
)

func main() {
	cmd.Execute()
}
