package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)


type Config struct {
	Base     BaseConfig     `mapstructure:"base"`
	Oss      OssConfig      `mapstructure:"oss"`
	BaseAuth BaseAuthConfig `mapstructure:"basic_auth"`
}

type BaseConfig struct {
	ListenAddress  string `mapstructure:"listen_address"`
	UploadDir      string `mapstructure:"upload_dir"`
	URLPrefix      string `mapstructure:"url_prefix"`
	DefaultStorage string `mapstructure:"default_storage"`
}

type OssConfig struct {
	Enable          bool   `mapstructure:"enable"`
	Public          bool   `mapstructure:"public"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key"`
	AccessKeySecret string `mapstructure:"access_secret"`
	BucketName      string `mapstructure:"bucket_name"`
}

type BaseAuthConfig struct {
	Enable   bool   `mapstructure:"enable"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (c *Config) validate() {
	if config.Base.UploadDir == "" && config.Base.DefaultStorage == storageLocal {
		logrus.Fatal("select local storage, but update path is null")
	}

	if !config.Oss.Enable && config.Base.DefaultStorage == storageAliyunOss {
		logrus.Fatal("select aliyun oss storage, but oss config is null ?")
	}
}

func (c *Config) print() {
	logrus.Infof("upload_dir:			%s", config.Base.UploadDir)
	logrus.Infof("baseauth.enable: 		%v", config.BaseAuth.Enable)
	logrus.Infof("default_storage:  		%s", config.Base.DefaultStorage)
	logrus.Infof("oss.enable:       		%v", config.Oss.Enable)
}

func parseConfig() *Config {
	logrus.Infof("config file: %s", configFile)

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc")
		viper.AddConfigPath("/etc/filestorage")
		viper.SetConfigName("filestorage")
	}

	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("can not read config file, err: %v", err)
	}
	logrus.Infof("read config from %s", viper.ConfigFileUsed())

	cfg := &Config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		logrus.Fatal("can not parse config file")
	}

	return cfg
}