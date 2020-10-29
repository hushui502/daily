#项目结构说明

- **tls-off:** 不开TLS验证的样例 (以连接调用“私密通讯”项目的链码为例)
  - fixtures 配置文件
    -  crypto-config-chat : “私密通讯”项目的crypto-config文件夹
    -  config-without-tls.yaml : sdk配置文件
  - main-chat.go 主函数
  
- **tls-on** 开启TLS验证的样例(以连接官方的样例项目为例)
  - fixtures 配置文件
    -  crypto-config : 官方样例项目的crypto-config文件夹
    -  config.yaml : sdk配置文件
  - main-e2e.go 主函数
  
- **keystore:** 程序运行过程中自动生成的文件夹，用于临时存放证书文件

- **go.mod** 项目依赖
  - **go.sum** 项目依赖