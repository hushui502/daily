package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var RuntimeViper *viper.Viper

func init() {
	RuntimeViper = viper.New()
	RuntimeViper.SetConfigType("toml")
	RuntimeViper.SetConfigName("cfg")
	RuntimeViper.AddConfigPath("etc/proxy")
	RuntimeViper.AddConfigPath("./config")
	if err := RuntimeViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file %s", err))
	}
	RuntimeViper.WatchConfig()
	RuntimeViper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config file changed %s", e.Name)
	})
}
