package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	configPath string
)

type Config struct {
	ListenAddr string `json:"listen"`
	RemoteAddr string `json:"remote"`
	Password   string `json:"password"`
}

func init() {
	home, _ := homedir.Dir()
	// 默认的配置文件夹名称
	configFilename := ".shadowsocks.json"
	// 如果有用户定义配置文件，则使用用户定义的
	if len(os.Args) == 2 {
		configFilename = os.Args[1]
	}

	configPath = path.Join(home, configFilename)
}

// 保存配置文件
func (config *Config) SaveConfig() {
	configJson, _ := json.MarshalIndent(config, "", "		")
	err := ioutil.WriteFile(configPath, configJson, 0644)
	if err != nil {
		fmt.Errorf("保存配置到文件 %s 出错: %s", configPath, err)
	}
	log.Printf("保存配置到文件 %s 成功\n", configPath)
}

func (config *Config) ReadConfig() {
	// 如果配置文件存在，就读取配置文件中的配置 assign 到 config
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		log.Printf("从文件 %s 中读取配置\n", configPath)
		file, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("打开配置文件 %s 出错:%s", configPath, err)
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(config)
		if err != nil {
			log.Fatalf("格式不合法的 JSON 配置文件:\n%s", file.Name())
		}
	}
}