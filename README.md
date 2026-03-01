# Orbit — A CLI tool for building go applications.

Orbit 是一个基于Golang的应用脚手架，它是由Golang生态中各种非常流行的库整合而成的，它可以帮助你快速构建一个高效、可靠的应用程序。


## 要求
要使用 Orbit，您需要在系统上安装以下软件：

* Golang 1.19或更高版本
* Git



### 安装

您可以通过以下命令安装 Orbit：

```bash
go install github.com/khy20040121/orbit@latest
```

国内用户可以使用`GOPROXY`加速`go install`

```
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
```

> tips: 如果`go install`成功，却提示找不到 orbit 命令，这是因为环境变量没有配置，可以把 GOBIN 目录配置到环境变量中即可


### 创建新项目

您可以使用以下命令创建一个新的Golang项目：

```bash
orbit new projectName
```


> Orbit 内置了两种类型的Layout：

* **基础模板(Basic Layout)**

Basic Layout 包含一个非常精简的架构目录结构 

* **高级模板(Advanced Layout)**

Advanced Layout 包含一个完善的架构目录结构


Advanced Layout 包含了很多 orbit 的用法示例（ db、redis、 jwt、 cron、 migration )等

此命令将创建一个名为`projectName`的目录，并在其中生成一个优雅的Golang项目结构。

### 创建组件

您可以使用以下命令为项目创建handler、service、repository和model等组件：

```bash
orbit create handler user
orbit create service user
orbit create repository user
orbit create model user
```
或
```
orbit create all user
```
这些命令将分别创建一个名为`UserHandler`、`UserService`、`UserRepository`和`UserModel`的组件，并将它们放置在正确的目录中。

### 启动项目

您可以使用以下命令快速启动项目：

```bash
orbit run
```

此命令将启动您的Golang项目。

### 编译wire.go

您可以使用以下命令快速编译`wire.go`：

```bash
orbit wire
```

此命令将编译您的`wire.go`文件，并生成所需的依赖项。
