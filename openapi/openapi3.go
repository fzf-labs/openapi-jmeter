package openapi

import (
	"context"
	"fmt"
	"log"

	"github.com/fzf-labs/openapi-jmeter/config"
	"github.com/fzf-labs/openapi-jmeter/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/sync/errgroup"
)

type OpenAPI3 struct {
	conf *config.Config
}

func NewOpenAPI3(conf *config.Config) *OpenAPI3 {
	return &OpenAPI3{
		conf: conf,
	}
}

func (o *OpenAPI3) Run() ([]*API, error) {
	// 从输入目录读取所有OpenAPI文件
	files, err := utils.ReadDirFilesWithSuffix(o.conf.Jmeter.InputPath, ".openapi.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI files: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no OpenAPI files found in %s", o.conf.Jmeter.InputPath)
	}

	apis := make([]*API, 0)
	// 创建errgroup用于并发处理
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(100) // 最大并发数100

	for _, file := range files {
		file := file
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				fileAPIs, err := o.processFile(file)
				if err != nil {
					log.Printf("Error processing file %s: %v", file, err)
					return nil
				}
				apis = append(apis, fileAPIs...)
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("error during OpenAPI files processing: %w", err)
	}

	log.Print("jmeter script generation completed")
	return apis, nil
}

func (o *OpenAPI3) processFile(file string) ([]*API, error) {
	// 读取OpenAPI文件
	openapiJSON, err := utils.ReadFileToString(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI file: %w", err)
	}

	// 解析OpenAPI规范
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData([]byte(openapiJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI JSON: %w", err)
	}
	utils.JsonDump(doc)
	apis := make([]*API, 0)
	return apis, nil
}
