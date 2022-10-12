package main

import (
	"github.com/saperliu/common-tool/common"
	"github.com/saperliu/common-tool/logger"
	"github.com/spf13/viper"
)

func main() {

}

func getConfig() *viper.Viper {
	//读取配置文件
	workPath := common.GetProgramPath("conf")
	logger.Info(" work path", workPath)
	serviceConfig := viper.New()
	serviceConfig.SetConfigName("config") // name of config file (without extension)
	serviceConfig.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	serviceConfig.AddConfigPath(workPath) // path to look for the config file in
	if err := serviceConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			logger.Error("config file not found")
			panic(err)
		} else {
			// Config file was found but another error was produced
			logger.Error("Config file was found but another error was produced")
			panic(err)
		}
	}
	return serviceConfig
}
