package config

import (
	"bytes"
	"fmt"
	"strconv"

	"nightingale-club203/src/common/address"
	"nightingale-club203/src/common/identity"
	"nightingale-club203/src/common/loggeri"
	"nightingale-club203/src/common/report"
	"nightingale-club203/src/modules/judge/backend/query"
	"nightingale-club203/src/modules/judge/backend/redi"
	"nightingale-club203/src/modules/judge/stra"

	"github.com/spf13/viper"
	"github.com/toolkits/pkg/file"
)

type ConfYaml struct {
	Logger            loggeri.Config           `yaml:"logger"`
	Query             query.SeriesQuerySection `yaml:"query"`
	Redis             redi.RedisSection        `yaml:"redis"`
	Strategy          stra.StrategySection     `yaml:"strategy"`
	Identity          identity.Identity        `yaml:"identity"`
	Report            report.ReportSection     `yaml:"report"`
	NodataConcurrency int                      `yaml:"nodataConcurrency"`
}

var (
	Config *ConfYaml
)

func Parse(conf string) error {
	bs, err := file.ReadBytes(conf)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(bs))
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	viper.SetDefault("query", map[string]interface{}{
		"maxConn":          100,
		"maxIdle":          10,
		"connTimeout":      1000,
		"callTimeout":      2000,
		"indexCallTimeout": 2000,
		"indexMod":         "index",
		"indexPath":        "/api/index/counter/clude",
	})

	viper.SetDefault("redis.idle", 5)
	viper.SetDefault("redis.prefix", "/n9e")
	viper.SetDefault("redis.timeout", map[string]int{
		"conn":  500,
		"read":  3000,
		"write": 3000,
	})

	viper.SetDefault("strategy", map[string]interface{}{
		"partitionApi":   "/api/mon/stras/effective?instance=%s:%s",
		"updateInterval": 9000,
		"indexInterval":  60000,
		"timeout":        5000,
		"mod":            "monapi",
		"eventPrefix":    "n9e",
	})

	viper.SetDefault("report", map[string]interface{}{
		"mod":      "judge",
		"enabled":  true,
		"interval": 4000,
		"timeout":  3000,
		"api":      "api/hbs/heartbeat",
		"remark":   "",
	})

	viper.SetDefault("nodataConcurrency", 1000)
	viper.SetDefault("pushUrl", "http://127.0.0.1:2058/v1/push")

	err = viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v\n", conf, err)
	}

	Config.Report.HTTPPort = strconv.Itoa(address.GetHTTPPort("judge"))
	Config.Report.RPCPort = strconv.Itoa(address.GetRPCPort("judge"))

	return err
}
