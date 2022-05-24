package main

import (
	"mixin_nft_cli/cmd"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("run nft cli ...")
	cmd.Run()
}
