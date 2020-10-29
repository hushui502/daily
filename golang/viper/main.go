package main

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)


func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	viper.WatchConfig()

	fmt.Println("redis port before sleep: ", viper.Get("redis.port"))
	time.Sleep(time.Second * 10)
	fmt.Println("redis port after sleep: ", viper.Get("redis.port"))

}
