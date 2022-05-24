package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

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

	root.PersistentFlags().StringP("config", "c", "", "configuration file path, default is ./config.yaml")
}

func initConfig() {
	logrus.Infoln("loading configuration ....")

	cp, err := root.Flags().GetString("config")
	if err != nil {
		panic(err)
	}

	if cp == "" {
		cp, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		logrus.Infoln("current dir:", cp)
		cp = filepath.Join(cp, "config.yaml")
		logrus.Infoln("config file path:", cp)
	}

	config.Load(cp, &cfg)
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
