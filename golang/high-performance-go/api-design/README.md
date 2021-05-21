# API设计

## API项目结构与管理

### API定义方式
使用grpc作为内部通信方式，因为protobuf即使代码也是文档，避免更新代码没有更新文档导致的交接对接困难。


![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/546bf12653ac440aa15b7abee10e1a23~tplv-k3u1fbpfcp-watermark.image)

### API Project

![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/07eb997e9f8f4a53ad75c52cc0715f34~tplv-k3u1fbpfcp-watermark.image)

工作流程
- 开发同事修改了proto文件定义之后push到对应的业务应用仓库中
- 触发cicd流程将proto文件复制到api project中
    - 首先会对proto文件进行静态代码分析，查看是否符合规范
    - 然后会clone api project创建一个新的分支
    - push代码，创建一个pr
- 然后我们对应负责的同学收到 code review 的通知之后进行 code review，没有问题就会合并到 api project 的主分支当中了
- 然后就会触发 cicd 生成对应语言的客户端代码，push 到对应的各个子仓库当中了

### API Project Layout


![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/bc83a6a3ebc14b42bb877163b6fe246f~tplv-k3u1fbpfcp-watermark.image)

- 首先是在业务项目当中，我们顶层会有一个 api 目录
    - 在 api 目录当中我们会按照 product name/app name/版本号/app.proto 的方式进行组织
    - 具体怎么组织可能每个公司都不太一样，但是总的来说就是应用的 唯一名称+版本号 来进行一个区分
- 在 api project 当中和业务应用类似，也有一个 api 目录，通过上图的两个框就可以发现这是一模一样的
    - 除此之外 api project 还有用于注解的 annotations 文件夹
    - 有一些第三方的引用，例如 googleapis 当中的一些 proto 文件

## API设计
### API兼容性设计

向下兼容的变更
- 新增接口
- 新增参数字段
- 新增返回字段

向下不兼容的变更（破坏性变更）
- 删除或重命名服务，字段，方法或枚举值
- 修改字段的类型
- 修改现有请求的可见行为
- 给资源消息添加 读取/写入 字段

### API 命名规范
包名

![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/f952cf4b796b46388369664f335fa30d~tplv-k3u1fbpfcp-watermark.image)

API 定义
- 方法+资源


![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/ef443ed0d21c4dec9d4f0ff24a5081f9~tplv-k3u1fbpfcp-watermark.image)


```js
// api/product/app/v1/blog.proto

syntax = "proto3";

package product.app.v1;

import "google/api/annotations.proto";

// blog service is a blog demo
service BlogService {

	rpc GetArticles(GetArticlesReq) returns (GetArticlesResp) {
		option (google.api.http) = {
			get: "/v1/articles"
			additional_bindings {
				get: "/v1/author/{author_id}/articles"
			}
		};
	}
}
```

一般而言我们应该为每个接口都创建一个自定义的 message，为了后面扩展，如果我们用 Empty 的话后续就没有办法新增字段了


### API Error

![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/bca9ed549b444525a6a42ae970f2148d~tplv-k3u1fbpfcp-watermark.image)


```js
message Status {
  // 错误码，跟 grpc-status 一致，并且在HTTP中可映射成 http-status
  int32 code = 1;
  // 错误原因，定义为业务判定错误码
  string reason = 2;
  // 错误信息，为用户可读的信息，可作为用户提示内容
  string message = 3;
  // 错误详细信息，可以附加自定义的信息列表
  repeated google.protobuf.Any details = 4;
}
```

## 错误传播
错误传播这一部分很容易出的问题就是，当前服务直接把上游服务的错误给返回了，这样会导致一些问题：
- 如果我调用了多个上游服务都报错了，我应该返回哪一个错误
- 直接返回导致必须要有一个全局错误码，不然的话就会冲突，但是全局错误码是很难定义的
  正确的做法应该是把上游错误信息吞掉，返回当前服务自己定义的错误信息就可以了。

## 实战-基于protobuf自动生成gin代码

gin example

```js
package main

import "github.com/gin-gonic/gin"

func handler(ctx *gin.Context) {
	// get params
	params := struct {
		Msg string `json:"msg"`
	}{}
	ctx.BindQuery(&params)

	// 业务逻辑

	// 返回数据
	ctx.JSON(200, gin.H{
		"message": params.Msg,
	})
}

func main() {
	r := gin.Default()
	r.GET("/ping", handler)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
```

grpc server interface

```js
// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}
```

```js
type GreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	mustEmbedUnimplementedGreeterServer()
}
```

### 方案
我们需要从 proto 文件中得知 http method，http path 的信息，这样我们才知道要注册到哪个路由上

大概的样子

```js
type GreeterService struct {
	server GreeterHTTPServer
	router gin.IRouter
}

// 生成的 gin.HandlerFunc
// 由于 HandlerFunc 签名的限制，不能从参数传递 service 接口进来
// 所以我们使用一个 Struct 托管 service 数据
func (s *GreeterService) SayHello(ctx *gin.Context) {
	var in HelloRequest

	if err := ctx.ShouldBindJSON(∈); err != nil {
		// 返回参数错误
		return
	}

	// 调用业务逻辑
	out, err := s.server.(GreeterHTTPServer).SayHello(ctx, ∈)
	if err != nil {
		// 返回错误结果
		return
	}

	// 返回成功结果
	ctx.JSON(200, out)
	return
}

// 路由注册，首先需要 gin.IRouter 接口用于注册
// 其次需要获取到 SayHello 方法对应的 http method 和 path
func (s *GreeterService) RegisterService() {
	s.router.Handle("GET", "/hello", s.SayHello)
}
```

实现

proto文件


```js
syntax = "proto3";

option go_package = "github.com/mohuishou/protoc-gen-go-gin/example/testproto;testproto";

package testproto;

import "google/api/annotations.proto";

// blog service is a blog demo
service BlogService {
	// 方法名 action+resource
	rpc GetArticles(GetArticlesReq) returns (GetArticlesResp) {
  	// 添加 option 用于指定 http 的路由和方法
		option (google.api.http) = {
			get: "/v1/articles"

      // 可以通过添加 additional_bindings 一个 rpc method 对应多个 http 路由
			additional_bindings {
				get: "/v1/author/{author_id}/articles"
			}
		};
	}
}
```

service info

```js
type service struct {
	Name     string // Greeter
	FullName string // helloworld.Greeter
	FilePath string // api/helloworld/helloworld.proto

	Methods   []*method
	MethodSet map[string]*method
}
```

method info

```js
type method struct {
	Name    string // SayHello
	Num     int    // 一个 rpc 方法可以对应多个 http 请求
	Request string // SayHelloReq
	Reply   string // SayHelloResp
	// http_rule
	Path         string // 路由
	Method       string // HTTP Method
	Body         string
	ResponseBody string
}
```

获取所有的 proto 文件

```js
// main.go
func main() {
	// ...

	options := protogen.Options{
		ParamFunc: flags.Set,
	}

	options.Run(func(gen *protogen.Plugin) error {
		// ...
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}
```

生成单个 proto 文件中的内容

```js
// 后面都是 gin.go 的内容

// generateFile generates a _gin.pb.go file.
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	// 如果不存在 service 就直接跳过了，我们主要生成 service 的接口
    if len(file.Services) == 0 {
		return nil
	}

	filename := file.GeneratedFilenamePrefix + "_gin.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by github.com/mohuishou/protoc-gen-go-gin. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the mohuishou/protoc-gen-go-gin package it is being compiled against.")
	g.P("// ", contextPkg.Ident(""), metadataPkg.Ident(""))
	g.P("//", ginPkg.Ident(""), errPkg.Ident(""))
	g.P()

	for _, service := range file.Services {
		genService(gen, file, g, service)
	}
	return g
}
```

获取 service 相关信息

```js
func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	if s.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	// HTTP Server.
	sd := &service{
		Name:     s.GoName,
		FullName: string(s.Desc.FullName()),
		FilePath: file.Desc.Path(),
	}

	for _, method := range s.Methods {
		sd.Methods = append(sd.Methods, genMethod(method)...)
	}
	g.P(sd.execute())
}
```

获取 rpc 方法的相关信息

```js
func genMethod(m *protogen.Method) []*method {
	var methods []*method

	// 存在 http rule 配置
	rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
	if rule != nil && ok {
		for _, bind := range rule.AdditionalBindings {
			methods = append(methods, buildHTTPRule(m, bind))
		}
		methods = append(methods, buildHTTPRule(m, rule))
		return methods
	}

	// 不存在走默认流程
	methods = append(methods, defaultMethod(m))
	return methods
}
```

从 option 中生成路由

```js
func buildHTTPRule(m *protogen.Method, rule *annotations.HttpRule) *method {
	// ....

	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		method = "GET"
	// ... 其他映射
	}
	md := buildMethodDesc(m, method, path)
	return md
}
```

生成默认的路由

```js
func defaultMethod(m *protogen.Method) *method {
    // 分割方法名
	names := strings.Split(toSnakeCase(m.GoName), "_")

    // ...

    // 如果 http method 映射成功，那么路由就是 names[1:]
    // 如果没有映射成功路由就是 names
	switch strings.ToUpper(names[0]) {
	case http.MethodGet, "FIND", "QUERY", "LIST", "SEARCH":
		httpMethod = http.MethodGet
	// ...  其他方法映射
	default:
		httpMethod = "POST"
		paths = names
	}

	// ...

	md := buildMethodDesc(m, httpMethod, path)
	return md
}
```

## ref
http://lailin.xyz/post/go-training-week4-config.html