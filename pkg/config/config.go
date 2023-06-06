package config

import (
	"errors"
	"github.com/spf13/viper"
)

func LoadConfig(config string) *viper.Viper {
	// var imagesList types.ImageList
	vip := viper.New()
	vip.SetConfigFile(config)
	if err := vip.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			panic(errors.New("没找到配置文件"))
		} else {
			// 配置文件被找到，但产生了另外的错误
			panic(err.Error())
		}
	}
	return vip
}
