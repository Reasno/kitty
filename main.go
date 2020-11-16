//go:generate trs proto/app.proto
package main

import (
	"glab.tagtic.cn/ad_gains/kitty/cmd"
)

func main() {
	cmd.Execute()
}
