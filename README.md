# FJmeter

FJmeter 是一个用Go语言编写的工具，用于自动生成JMeter测试脚本。它可以帮助测试人员快速创建性能测试用例，特别适用于API压力测试场景。

## 功能特点

- 自动生成JMeter测试脚本
- 支持自定义线程数、启动时间和持续时间
- 支持通过关键字筛选需要生成测试脚本的API
- 灵活的输入输出路径配置
- 支持设置域名和用户信息

## 安装

确保你已经安装了Go 1.19或更高版本，然后运行：

```bash
go get github.com/fzf-labs/fjmeter
```

## 使用方法

### 命令行参数

```bash
fjmeter [flags]
```

### 可用的参数：

- `-d, --domain string`：设置域名
- `-u, --user string`：设置用户信息
- `-n, --numThreads int`：设置线程数（默认：200）
- `-r, --rampTime int`：设置线程启动时间（默认：300秒）
- `-c, --continueTime int`：设置持续时间（默认：600秒）
- `-i, --inPutPath string`：设置输入路径（默认："./api"）
- `-o, --outPutPath string`：设置输出路径（默认："./doc/jmeter"）
- `-k, --keyword string`：设置匹配关键字（默认："压测"）

### 示例

生成基本的JMeter测试脚本：
```bash
fjmeter -d "example.com" -u "testuser" -n 100 -r 60 -c 300
```

使用自定义路径：
```bash
fjmeter -i "./myapi" -o "./output/jmeter" -k "performance"
```

## 依赖

- github.com/go-openapi/spec v0.21.0
- github.com/samber/lo v1.47.0
- github.com/spf13/cobra v1.8.1

## 项目结构

```
fjmeter/
├── main.go          # 主程序入口
├── jmeter/          # JMeter相关功能实现
├── tpl/            # 模板文件目录
├── go.mod          # Go模块定义
└── go.sum          # 依赖版本锁定文件
```

## 许可证

[License信息]

## 贡献

欢迎提交问题和Pull Request！

## 作者

[作者信息]
