# OpenAPI-JMeter

OpenAPI-JMeter 是一个基于 Go 语言开发的工具，用于自动将 OpenAPI（Swagger）文档转换为 JMeter 测试脚本。它能帮助测试人员快速创建性能测试用例，特别适用于 API 压力测试场景。

## 功能特点

- 支持 OpenAPI 2.0 和 3.0 规范
- 自动生成 JMeter 测试脚本（.jmx 文件）
- 支持自定义 HTTP 请求配置（协议、服务器、端口等）
- 灵活的线程组配置（线程数、启动时间、持续时间等）
- 支持 CSV 数据集配置
- 内置后端监听器支持（InfluxDB、Graphite）
- 支持查看结果树配置
- 支持通过关键字筛选需要生成测试脚本的 API
- 灵活的输入输出路径配置
- 支持并发处理多个OpenAPI文件
- 自动生成CSV数据集文件
- 支持查看测试结果

## 安装

确保你已经安装了 Go 1.20 或更高版本，然后运行：

```
go install github.com/fzf-labs/openapi-jmeter
```

## 配置文件

默认配置文件路径为 `./config/config.yaml`：

### 基础配置
```yaml
jmeter:
  openapiVersion: 3.0  # 2.0 表示 OpenAPI 2.0，3.0 表示 OpenAPI 3.0
  inputPath: "./example"  # OpenAPI 文件路径
  outputPath: "./example/jmx" # JMeter 输出文件路
  outputMode: "skip" # 文件输出模式，overwrite 表示覆盖，skip 表示跳过
  keyword: ""  # 需要生成测试用例的接口名称关键字 留空表示不过滤，全部生成。
  suffix: ".openapi.yaml" # OpenAPI 文件后缀
# Jmeter 基础配置
# HTTP 请求配置
httpRequest:
  protocol: "http"  # http, https
  serverNameOrIp: "localhost" # 服务器名称或IP
  portNumber: "8080" # 端口号
  redirectAutomatically: false # 自动重定向
  followRedirects: true # 跟随重定向
  useKeepAlive: false # 使用保持连接
  useMultipartFormData: false # 使用多部分表单数据
  browserCompatibleHeaders: false # 浏览器兼容头
  httpDefaultHeaders: # HTTP 默认请求头
    - key: "Content-Type"
      value: "application/json"
    - key: "Accept"
      value: "application/json"
    - key: "User-Agent" 
      value: "JMeter"
# 线程组配置
threadGroup:
  actionToBeTakenAfterASamplerError: "continue"  # continue, startnextloop, stopthread, stoptest
  numThreads: 100                                 # 线程数量
  rampTime: 20                                   # 线程启动时间(秒)
  loopCount: -1                                  # -1 表示永远循环
  sameUserOnEachIteration: false                 # 每次迭代使用相同用户
  delayThreadCreationUntilNeeded: false          # 延迟线程创建直到需要
  specifyThreadLifetime: true                    # 指定线程生命周期
  duration: 300                                  # 持续时间(秒)
  startupDelay: 0                                # 启动延迟(秒)
# CSV 数据集配置
csvDataSetConfig:
  fileNamePrefix: "./example" # 文件路径前缀
  fileEncoding: "UTF-8" # 文件编码格式
  ignoreFirstLine: true # 是否忽略第一行
  delimiter: "," # 分隔符
  allowQuotedData: false # 是否允许带引号的数据
  recycle: true # 是否循环使用数据
  stopThread: false # 数据用完时是否停止线程
  shareMode: "shareMode.all"  # shareMode.all, shareMode.group, shareMode.thread
# 后端监听器配置
backendListener:
  enable: true # 是否启用后端监听器
  backendListenerImplementation: "org.apache.jmeter.visualizers.backend.influxdb.InfluxdbBackendListenerClient" # 后端监听器实现
  asyncQueueSize: 5000 # 异步队列大小 
  # Graphite 配置
  graphite:
    graphiteMetricsSender: "org.apache.jmeter.visualizers.backend.graphite.TextGraphiteMetricsSender" # 图表配置  
    graphiteHost: "localhost" # 图表主机
    graphitePort: "2003" # 图表端口
    rootMetricsPrefix: "jmeter." # 根度量前缀
    summaryOnly: false # 仅总结
    samplersList: ".*" # 采样器列表
    useRegexpForSamplersList: true # 使用正则表达式
    percentiles: "90;95;99" # 百分位数
  # InfluxDB Raw 配置
  influxdbRaw:
    influxdbMetricsSender: "org.apache.jmeter.visualizers.backend.influxdb.HttpMetricsSender" # InfluxDB 配置
    influxdbUrl: "http://localhost:8086/write?db=jmeter" # InfluxDB URL
    influxdbToken: "your-token-here" # InfluxDB 令牌
    measurement: "jmeter" # 度量
  # InfluxDB 配置
  influxdb:
    influxdbMetricsSender: "org.apache.jmeter.visualizers.backend.influxdb.HttpMetricsSender" # InfluxDB 配置
    influxdbUrl: "http://localhost:8086/api/v2/write" # InfluxDB URL
    application: "MyTestApp" # 应用
    measurement: "jmeter" # 度量
    summaryOnly: "false" # 仅总结
    samplersRegex: ".*" # 采样器正则
    percentiles: "90;95;99" # 百分位数
    testTitle: "API Performance Test" # 测试标题
    eventTags: "" # 事件标签
# 查看结果树配置
viewResultsTree:
  enable: true # 是否启用查看结果树
  fileNamePrefix: "./example" # 文件路径前缀
  logDisplayOnly: "false" # 仅显示日志  
```
## 使用方法

### 命令行运行

```bash
openapi-jmeter -c ./config/config.yaml
```

### 命令行参数

- `-c, --config string`：配置文件路径（默认："./config/config.yaml"）

## 项目结构

```
openapi-jmeter/
├── config/         # 配置管理
│   ├── config.go   # 配置定义和解析
│   └── config.yaml # 默认配置文件
├── jmeter/         # JMeter脚本生成
│   └── jmeter.go   # 核心生成逻辑
├── openapi/        # OpenAPI解析
│   ├── openapi.go  # 版本选择
│   ├── openapi2.go # 2.0规范支持
│   └── openapi3.go # 3.0规范支持  
├── utils/          # 工具函数
├── tpl/           # JMeter模板
├── main.go        # 程序入口
└── go.mod         # 依赖管理
```

## 贡献

欢迎提交问题和 Pull Request！

## 许可证

MIT License