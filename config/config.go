package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Jmeter           Jmeter           `yaml:"jmeter" json:"jmeter"`                     // JMeter 配置
	HttpRequest      HttpRequest      `yaml:"httpRequest" json:"httpRequest"`           // HTTP 请求配置
	ThreadGroup      ThreadGroup      `yaml:"threadGroup" json:"threadGroup"`           // 线程组配置
	CsvDataSetConfig CsvDataSetConfig `yaml:"csvDataSetConfig" json:"csvDataSetConfig"` // CSV 数据集配置
	BackendListener  BackendListener  `yaml:"backendListener" json:"backendListener"`   // 后端监听器配置
	ViewResultsTree  ViewResultsTree  `yaml:"viewResultsTree" json:"viewResultsTree"`   // 结果树配置
}

type Jmeter struct {
	OpenapiVersion string `yaml:"openapiVersion" json:"openapiVersion"` // OpenAPI 版本
	InputPath      string `yaml:"inputPath" json:"inputPath"`           // OpenAPI 路径
	OutputPath     string `yaml:"outputPath" json:"outputPath"`         // JMeter 输出文件路径
	OutputMode     string `yaml:"outputMode" json:"outputMode"`         // 文件输出模式，overwrite 表示覆盖，skip 表示跳过
	Keyword        string `yaml:"keyword" json:"keyword"`               // 关键字过滤
	Suffix         string `yaml:"suffix" json:"suffix"`                 // OpenAPI 文件后缀
}

type HttpRequest struct {
	Protocol                 string            `yaml:"protocol" json:"protocol"`                                 // 协议
	ServerNameOrIp           string            `yaml:"serverNameOrIp" json:"serverNameOrIp"`                     // 服务器名称或IP
	PortNumber               string            `yaml:"portNumber" json:"portNumber"`                             // 端口号
	RedirectAutomatically    bool              `yaml:"redirectAutomatically" json:"redirectAutomatically"`       // 自动重定向
	FollowRedirects          bool              `yaml:"followRedirects" json:"followRedirects"`                   // 跟随重定向
	UseKeepAlive             bool              `yaml:"useKeepAlive" json:"useKeepAlive"`                         // 使用保持连接
	UseMultipartFormData     bool              `yaml:"useMultipartFormData" json:"useMultipartFormData"`         // 使用多部分表单数据
	BrowserCompatibleHeaders bool              `yaml:"browserCompatibleHeaders" json:"browserCompatibleHeaders"` // 浏览器兼容头
	HTTPDefaultHeaders       []HTTPHeadersItem `yaml:"httpDefaultHeaders" json:"httpDefaultHeaders"`             // http默认请求头
}

type HTTPHeadersItem struct {
	Key   string `yaml:"key" json:"key"`
	Value string `yaml:"value" json:"value"`
}

type ThreadGroup struct {
	ActionToBeTakenAfterASamplerError string `yaml:"actionToBeTakenAfterASamplerError" json:"actionToBeTakenAfterASamplerError"` // 采样器错误后的操作
	NumThreads                        int    `yaml:"numThreads" json:"numThreads"`                                               // 线程数量
	RampTime                          int    `yaml:"rampTime" json:"rampTime"`                                                   // 线程启动时间(秒)
	LoopCount                         int    `yaml:"loopCount" json:"loopCount"`                                                 // 循环次数
	SameUserOnEachIteration           bool   `yaml:"sameUserOnEachIteration" json:"sameUserOnEachIteration"`                     // 每次迭代使用相同用户
	DelayThreadCreationUntilNeeded    bool   `yaml:"delayThreadCreationUntilNeeded" json:"delayThreadCreationUntilNeeded"`       // 延迟线程创建直到需要
	SpecifyThreadLifetime             bool   `yaml:"specifyThreadLifetime" json:"specifyThreadLifetime"`                         // 指定线程生命周期
	Duration                          int    `yaml:"duration" json:"duration"`                                                   // 线程持续时间(秒)
	StartupDelay                      int    `yaml:"startupDelay" json:"startupDelay"`                                           // 线程启动延迟(秒)
}

type CsvDataSetConfig struct {
	FileNamePrefix  string `yaml:"fileNamePrefix" json:"fileNamePrefix"`   // 文件路径前缀
	FileEncoding    string `yaml:"fileEncoding" json:"fileEncoding"`       // 文件编码格式
	IgnoreFirstLine bool   `yaml:"ignoreFirstLine" json:"ignoreFirstLine"` // 是否忽略第一行
	Delimiter       string `yaml:"delimiter" json:"delimiter"`             // 分隔符
	AllowQuotedData bool   `yaml:"allowQuotedData" json:"allowQuotedData"` // 是否允许带引号的数据
	Recycle         bool   `yaml:"recycle" json:"recycle"`                 // 是否循环使用数据
	StopThread      bool   `yaml:"stopThread" json:"stopThread"`           // 数据用完时是否停止线程
	ShareMode       string `yaml:"shareMode" json:"shareMode"`             // 数据共享模式
}

type BackendListener struct {
	Enable                        bool        `yaml:"enable" json:"enable"`                                               // 是否启用后端监听器
	BackendListenerImplementation string      `yaml:"backendListenerImplementation" json:"backendListenerImplementation"` // 后端监听器实现类
	AsyncQueueSize                int         `yaml:"asyncQueueSize" json:"asyncQueueSize"`                               // 异步队列大小
	Graphite                      Graphite    `yaml:"graphite" json:"graphite"`                                           // Graphite 配置
	InfluxdbRaw                   InfluxdbRaw `yaml:"influxdbRaw" json:"influxdbRaw"`                                     // InfluxDB Raw 配置
	Influxdb                      Influxdb    `yaml:"influxdb" json:"influxdb"`                                           // InfluxDB 配置
}

type Graphite struct {
	GraphiteMetricsSender    string `yaml:"graphiteMetricsSender" json:"graphiteMetricsSender"`       // Graphite 指标发送器
	GraphiteHost             string `yaml:"graphiteHost" json:"graphiteHost"`                         // Graphite 服务器主机地址
	GraphitePort             string `yaml:"graphitePort" json:"graphitePort"`                         // Graphite 服务器端口
	RootMetricsPrefix        string `yaml:"rootMetricsPrefix" json:"rootMetricsPrefix"`               // 指标前缀
	SummaryOnly              bool   `yaml:"summaryOnly" json:"summaryOnly"`                           // 是否只发送汇总数据
	SamplersList             string `yaml:"samplersList" json:"samplersList"`                         // 采样器列表
	UseRegexpForSamplersList bool   `yaml:"useRegexpForSamplersList" json:"useRegexpForSamplersList"` // 是否使用正则表达式匹配采样器列表
	Percentiles              string `yaml:"percentiles" json:"percentiles"`                           // 百分比
}

type InfluxdbRaw struct {
	InfluxdbMetricsSender string `yaml:"influxdbMetricsSender" json:"influxdbMetricsSender"` // InfluxDB 指标发送器
	InfluxdbUrl           string `yaml:"influxdbUrl" json:"influxdbUrl"`                     // InfluxDB 服务器URL
	InfluxdbToken         string `yaml:"influxdbToken" json:"influxdbToken"`                 // InfluxDB 认证令牌
	Measurement           string `yaml:"measurement" json:"measurement"`                     // 测量指标名称
}

type Influxdb struct {
	InfluxdbMetricsSender string `yaml:"influxdbMetricsSender" json:"influxdbMetricsSender"` // InfluxDB 指标发送器
	InfluxdbUrl           string `yaml:"influxdbUrl" json:"influxdbUrl"`                     // InfluxDB 服务器URL
	Application           string `yaml:"application" json:"application"`                     // 应用名称
	Measurement           string `yaml:"measurement" json:"measurement"`                     // 测量指标名称
	SummaryOnly           string `yaml:"summaryOnly" json:"summaryOnly"`                     // 是否只发送汇总数据
	SamplersRegex         string `yaml:"samplersRegex" json:"samplersRegex"`                 // 采样器正则表达式
	Percentiles           string `yaml:"percentiles" json:"percentiles"`                     // 百分比
	TestTitle             string `yaml:"testTitle" json:"testTitle"`                         // 测试标题
	EventTags             string `yaml:"eventTags" json:"eventTags"`                         // 事件标签
}

type ViewResultsTree struct {
	Enable         bool   `yaml:"enable" json:"enable"`                 // 是否启用查看结果树
	FileNamePrefix string `yaml:"fileNamePrefix" json:"fileNamePrefix"` // 文件路径前缀
	LogDisplayOnly string `yaml:"logDisplayOnly" json:"logDisplayOnly"` // 仅显示日志
}

func NewConfig(file string) (*Config, error) {
	// 验证文件是否存在
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件 %s 不存在", file)
	}

	// 读取配置文件
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}
	// 解析配置
	var config Config
	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}
	return &config, nil
}
