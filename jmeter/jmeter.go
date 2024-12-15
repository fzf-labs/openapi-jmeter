package jmeter

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/fzf-labs/openapi-jmeter/config"
	"github.com/fzf-labs/openapi-jmeter/openapi"
	"github.com/fzf-labs/openapi-jmeter/tpl"
	"github.com/fzf-labs/openapi-jmeter/utils"
	"golang.org/x/sync/errgroup"
)

// JMeter JMeter配置结构体
type JMeter struct {
	Config *config.Config
	APIs   []*openapi.API
}

// JMXTemplateData 用于渲染JMX模板的数据结构
type JMXTemplateData struct {
	Config *config.Config
	API    *openapi.API
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

	// 解析模板
	tmpl, err := template.New("jmx").Parse(tpl.TplJmx)
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	// 创建errgroup用于并发处理
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(10) // 最大并发数10

	// 用于保护文件写入
	var mu sync.Mutex

	// 并发处理每个API
	for _, api := range j.APIs {
		api := api // 创建副本避免闭包问题
		g.Go(func() error {
			// 检查上下文是否已取消
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// 准备模板数据
			data := &JMXTemplateData{
				Config: j.Config,
				API:    api,
			}

			// 渲染模板
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				return fmt.Errorf("渲染模板失败: %v", err)
			}

			// 生成文件名
			fileName := fmt.Sprintf("%s_%s.jmx", api.HTTPDomain, api.HTTPName)
			filePath := filepath.Join(j.Config.Jmeter.OutputPath, fileName)

			// 加锁保护文件写入
			mu.Lock()
			err := utils.WriteContentCover(filePath, buf.String())
			mu.Unlock()

			if err != nil {
				return fmt.Errorf("写入JMX文件失败: %v", err)
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
