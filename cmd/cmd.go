package cmd

import (
	"encoding/json"
	"os"

	"mixin_nft_cli/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfg config.Config
var root cobra.Command = cobra.Command{
	Use: "mixin-nft-cli",
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	logrus.Infoln("loading config ....")
	config.Load("config.yaml", &cfg)
}

func prettyPrint(body interface{}) error {
	switch d := body.(type) {
	case []byte:
		var data map[string]interface{}
		if err := json.Unmarshal(d, &data); err != nil {
			return err
		}

		bts, err := json.MarshalIndent(data, " ", "    ")
		if err != nil {
			return err
		}

		logrus.Infof("%s\n", string(bts))
	default:
		bts, err := json.MarshalIndent(d, " ", "    ")
		if err != nil {
			return err
		}

		logrus.Infof("%s\n", string(bts))
	}

	return nil
}

func Run() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
