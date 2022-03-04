# IAM

五色神牛权限系统接口层，为了降低前后端交付成本的 API 框架, 支持服务抽象，命令行生成多语言 SDK。

## 特性

* 支持 Service 映射路由接口；
* 方便实现 Mock 接口；
* 集成 API 调试工具，比 Swagger 更好用；
* 支持自动生成文档；
* 支持自定义 SDK 生成器；  
* 支持通过命令行生成多语言 SDK。

## 包安装

```
go get -u github.com/utilslab/iam
```

## 命令行工具安装

```
go get -u github.com/utilslab/iam/cmd/iam
```

## 快速上手

参考 Foo 用例，[查看源码](https://github.com/koyeo/buck/tree/master/example/foo)。



## SDK 生成
以本例文档服务监听的 9090 端口为例。

**生成 Go SDK:**

```
$ iam sdk --address 127.0.0.0:9090 --output ./sdk  --package test-sdk --target go -y
```

**生成 Angular SDK:**

```
$ iam sdk --address 127.0.0.0:9090 --output ./sdk  --package test-sdk --target angular -y
```

**生成 Umi SDK:**

```
$ iam sdk --address 127.0.0.0:9090 --output ./sdk  --package test-sdk --target umi -y
```

## 服务方法

**格式说明:**

```
func (ctx context.Context [, in struct])([out interface{},] err error)
```

**入参说明:**

第一个参数必选，必须是 `context.Context` 类型，第二个参数可选，必须是一个 `struct`。

**出参说明:**

第一个参数 `out` 可选，可以为任意类型。第二个参数必选，必须是 error 类型。

**出参编码：**

如果 `out` 为 go 的基础数据类型，如 `string、int、float64` 等、或实现了 `String() `方法，则 API 返回报文采用字符串编码，否则将采用 json 编码并输出。


## 入参标签

| 标签        | 用途                                   |
| --------- | ------------------------------------ |
| label     | 用于备注字段在文档中的显示名称                      |
| validator | 用于标注字段的校验规则，如 `validator="required"` |

## Context Wrapper

通过 buck 实例调用 SetContextWrapper 方法，可以为引擎注入一个服务的 Context 包装器，以获得服务需要的上下文，如登录状态等。

```go
api.SetContextWrapper(function(c *gin.Context)(context.Context, error){
	return context.Background(), nil
})
```

## Authors 关于作者

- [**koyeo**](https://github.com/koeyo) - *Initial work*

查看更多关于这个项目的贡献者，请阅读 [contributors](https://gist.github.com/wangyan/6e8021667fe7f2082d153bed2d764618#)。

## License 授权协议

这个项目 MIT 协议， 请点击 [LICENSE](https://choosealicense.com/licenses/mit) 了解更多细节。
