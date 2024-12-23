package jmeter

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fzf-labs/openapi-jmeter/config"
	"github.com/fzf-labs/openapi-jmeter/openapi"
	"github.com/fzf-labs/openapi-jmeter/tpl"
	"github.com/fzf-labs/openapi-jmeter/utils"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

// JMeter JMeter配置结构体
type JMeter struct {
	Config *config.Config
	APIs   []*openapi.API
}

// JMXTemplateData 用于渲染JMX模板的数据结构
type JMXTemplateData struct {
	Config                    *config.Config
	API                       *openapi.API
	JmxFileName               string
	CSVFileName               string
	ViewResultsTreeFileName   string
	CSVVariableNamesParams    []string
	CSVVariableNamesParamsStr string
}

// NewJMeter 创建新的JMeter实例
func NewJMeter(cfg *config.Config, apis []*openapi.API) *JMeter {
	return &JMeter{
		Config: cfg,
		APIs:   apis,
	}
}

// GenerateJMX 生成JMX文件
func (j *JMeter) GenerateJMX() error {
	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(j.Config.Jmeter.OutputPath), 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}
	// 创建errgroup用于并发处理
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(10) // 最大并发数10
	// 并发处理每个API
	for _, v := range j.APIs {
		api := v // 创建副本避免闭包问题
		g.Go(func() error {
			// 检查上下文是否已取消
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			// 收集所有参数
			csvVariableNamesParams := make([]string, 0)
			// 添加header参数
			for _, header := range api.HTTPHeaders {
				csvVariableNamesParams = append(csvVariableNamesParams, header.Key)
			}
			// 添加query参数
			for _, param := range api.HTTPParams.Query {
				csvVariableNamesParams = append(csvVariableNamesParams, param.Key)
			}
			// 添加path参数
			for _, param := range api.HTTPParams.Path {
				csvVariableNamesParams = append(csvVariableNamesParams, param.Key)
			}
			// 添加JSON body参数
			if api.HTTPBody.JSON != nil {
				for _, param := range api.HTTPBody.JSON.Params {
					csvVariableNamesParams = append(csvVariableNamesParams, param.Key)
				}
			}
			// 参数去重
			csvVariableNamesParams = lo.Uniq(csvVariableNamesParams)

			data := &JMXTemplateData{
				Config:                    j.Config,
				API:                       api,
				JmxFileName:               strings.ToLower(fmt.Sprintf("%s_%s.jmx", api.HTTPDomain, api.HTTPName)), // 文件名称小写
				CSVFileName:               strings.ToLower(fmt.Sprintf("%s_%s.csv", api.HTTPDomain, api.HTTPName)),
				ViewResultsTreeFileName:   strings.ToLower(fmt.Sprintf("%s_%s.txt", api.HTTPDomain, api.HTTPName)),
				CSVVariableNamesParams:    csvVariableNamesParams,
				CSVVariableNamesParamsStr: strings.Join(csvVariableNamesParams, ","),
			}

			// 生成jmx文件
			if err := j.generateJMX(data); err != nil {
				return fmt.Errorf("生成JMX文件失败: %v", err)
			}
			// 生成csv文件
			if err := j.generateCSV(data); err != nil {
				return fmt.Errorf("生成CSV文件失败: %v", err)
			}
			// 生成viewResultsTree文件
			if err := j.generateViewResultsTree(data); err != nil {
				return fmt.Errorf("生成ViewResultsTree文件失败: %v", err)
			}
			return nil
		})
	}
	// 等待所有协程完成并处理错误
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

// generateJMX 生成jmx文件
func (j *JMeter) generateJMX(data *JMXTemplateData) error {
	// 生成jmx文件名
	filePath := filepath.Join(j.Config.Jmeter.OutputPath, data.JmxFileName)
	// 渲染模板
	var buf bytes.Buffer
	tmpl, err := template.New("jmx").Parse(tpl.TplJmx)
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染模板失败: %v", err)
	}
	// 加锁保护文件写入
	err = utils.WriteContentCover(filePath, buf.String())
	if err != nil {
		return fmt.Errorf("写入JMX文件失败: %v", err)
	}
	return nil
}

// generateCSV 生成csv文件
func (j *JMeter) generateCSV(data *JMXTemplateData) error {
	// 如果没有参数则不生成CSV
	if len(data.CSVVariableNamesParams) == 0 {
		return nil
	}
	// 生成csv文件名
	csvFilePath := filepath.Join(j.Config.Jmeter.OutputPath, data.CSVFileName)
	// 生成CSV内容
	var csvContent strings.Builder
	// 写入表头
	csvContent.WriteString(strings.Join(data.CSVVariableNamesParams, ","))
	csvContent.WriteString("\n")
	// 写入文件
	err := utils.WriteContentCover(csvFilePath, csvContent.String())
	if err != nil {
		return fmt.Errorf("写入CSV文件失败: %v", err)
	}
	return nil
}

// generateViewResultsTree 生成viewResultsTree文件
func (j *JMeter) generateViewResultsTree(data *JMXTemplateData) error {
	if !j.Config.ViewResultsTree.Enable {
		return nil
	}
	// 生成viewResultsTree文件名
	filePath := filepath.Join(j.Config.Jmeter.OutputPath, data.ViewResultsTreeFileName)
	// 生成viewResultsTree文件
	err := utils.WriteContentCover(filePath, "")
	if err != nil {
		return fmt.Errorf("创建ViewResultsTree文件失败: %v", err)
	}
	return nil
}
