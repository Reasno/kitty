//go:generate trs proto/app.proto
//go:generate trs proto/share.proto
package main

import (
	"glab.tagtic.cn/ad_gains/kitty/cmd"
)

func main() {
	cmd.Execute()
}
