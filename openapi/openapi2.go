package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/fzf-labs/openapi-jmeter/config"
	"github.com/fzf-labs/openapi-jmeter/utils"
	"github.com/getkin/kin-openapi/openapi2"
	"golang.org/x/sync/errgroup"
)

type OpenAPI2 struct {
	conf *config.Config
}

func NewOpenAPI2(conf *config.Config) *OpenAPI2 {
	return &OpenAPI2{
		conf: conf,
	}
}

func (o *OpenAPI2) Run() ([]*API, error) {
	// 从输入目录读取所有OpenAPI文件
	files, err := utils.ReadDirFilesWithSuffix(o.conf.Jmeter.InputPath, ".swagger.json")
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

func (o *OpenAPI2) processFile(file string) ([]*API, error) {
	// 读取OpenAPI文件
	swaggerJSON, err := utils.ReadFileToString(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI file: %w", err)
	}

	// 解析OpenAPI规范
	swagger := &openapi2.T{}
	if err := json.Unmarshal([]byte(swaggerJSON), swagger); err != nil {
		return nil, fmt.Errorf("failed to parse Swagger JSON: %w", err)
	}

	apis := make([]*API, 0)

	// 遍历所有路径
	for path, pathItem := range swagger.Paths {
		// 处理每种HTTP方法
		operations := map[string]*openapi2.Operation{
			"GET":    pathItem.Get,
			"POST":   pathItem.Post,
			"PUT":    pathItem.Put,
			"DELETE": pathItem.Delete,
		}

		for method, operation := range operations {
			if operation == nil {
				continue
			}

			fileName := o.cleanFileName(operation.Summary)
			// 检查关键字过滤
			if o.conf.Jmeter.Keyword != "" && !strings.Contains(fileName, o.conf.Jmeter.Keyword) {
				continue
			}

			api := &API{
				HTTPName:    fileName,
				HTTPDomain:  o.conf.HttpRequest.ServerNameOrIp,
				HTTPPath:    path,
				HTTPMethod:  method,
				HTTPHeaders: o.extractHeaders(operation),
				HTTPParams:  o.extractParams(operation),
				HTTPBody:    o.extractBody(swagger, operation),
			}
			apis = append(apis, api)
		}
	}

	return apis, nil
}

func (o *OpenAPI2) cleanFileName(filename string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9._\-\p{L}]`)
	return re.ReplaceAllString(filename, "")
}

func (o *OpenAPI2) extractHeaders(operation *openapi2.Operation) []*HTTPKeyAndType {
	headers := make([]*HTTPKeyAndType, 0)
	for _, param := range operation.Parameters {
		if param.In == "header" {
			headers = append(headers, &HTTPKeyAndType{
				Key:   param.Name,
				Value: "${" + param.Name + "}",
				Type:  param.Type.Slice()[0],
			})
		}
	}
	return headers
}

func (o *OpenAPI2) extractParams(operation *openapi2.Operation) HTTPParams {
	params := HTTPParams{
		Query: make([]*HTTPKeyAndType, 0),
		Path:  make([]*HTTPKeyAndType, 0),
	}

	for _, param := range operation.Parameters {
		switch param.In {
		case "query":
			params.Query = append(params.Query, &HTTPKeyAndType{
				Key:   param.Name,
				Value: "${" + param.Name + "}",
				Type:  param.Type.Slice()[0],
			})
		case "path":
			params.Path = append(params.Path, &HTTPKeyAndType{
				Key:   param.Name,
				Value: "${" + param.Name + "}",
				Type:  param.Type.Slice()[0],
			})
		}
	}

	return params
}

func (o *OpenAPI2) extractBody(swagger *openapi2.T, operation *openapi2.Operation) HTTPBody {
	body := HTTPBody{
		ContentType: "application/json",
		FormData:    make([]*HTTPKeyAndType, 0),
		JSON: &HTTPBodyJSON{
			JSONStr: "",
			Params:  make([]*HTTPKeyAndType, 0),
		},
	}

	for _, param := range operation.Parameters {
		switch param.In {
		case "formData":
			o.processFormData(param, &body)

		case "body":
			if err := o.processBodyParam(swagger, param, &body); err != nil {
				log.Printf("Process body param failed: %v", err)
			}
		}
	}

	return body
}

func (o *OpenAPI2) processFormData(param *openapi2.Parameter, body *HTTPBody) {
	body.ContentType = "multipart/form-data"
	body.FormData = append(body.FormData, &HTTPKeyAndType{
		Key:   param.Name,
		Value: "${" + param.Name + "}",
		Type:  param.Type.Slice()[0],
	})
}

func (o *OpenAPI2) processBodyParam(swagger *openapi2.T, param *openapi2.Parameter, body *HTTPBody) error {
	if param == nil || param.Schema == nil {
		return nil
	}

	// 需要处理 param *openapi2.Parameter 中 param.Type为 array 和 object 的情况

	jsonValue, err := o.processSchema(swagger, param.Schema, "", body)
	if err != nil {
		return fmt.Errorf("failed to process schema: %w", err)
	}

	jsonBytes, err := json.Marshal(jsonValue)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	body.JSON.JSONStr = string(jsonBytes)

	return nil
}

func (o *OpenAPI2) processSchema(swagger *openapi2.T, schema *openapi2.SchemaRef, parentPath string, body *HTTPBody) (interface{}, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	// 优先处理引用类型
	if schema.Ref != "" {
		parts := strings.Split(strings.TrimPrefix(schema.Ref, "#/"), "/")
		if len(parts) != 2 || parts[0] != "definitions" {
			return nil, fmt.Errorf("invalid reference format: %s", schema.Ref)
		}

		refSchema, ok := swagger.Definitions[parts[1]]
		if !ok {
			return nil, fmt.Errorf("reference not found: %s", parts[1])
		}

		if refSchema == nil || refSchema.Value == nil {
			return nil, fmt.Errorf("referenced schema is nil: %s", parts[1])
		}

		// 递归处理引用的schema
		return o.processSchema(swagger, refSchema, parentPath, body)
	}

	// 如果schema.Value为空,返回空值
	if schema.Value == nil {
		return nil, nil
	}

	// 获取类型信息
	schemaType := "string"
	if schema.Value.Type != nil && len(schema.Value.Type.Slice()) > 0 {
		schemaType = schema.Value.Type.Slice()[0]
	}

	// 根据类型处理
	switch schemaType {
	case "array":
		if schema.Value.Items == nil {
			return []interface{}{}, nil
		}
		itemValue, err := o.processSchema(swagger, schema.Value.Items, parentPath, body)
		if err != nil {
			return nil, fmt.Errorf("failed to process array item: %w", err)
		}
		return []interface{}{itemValue}, nil

	case "object":
		objMap := make(map[string]interface{})
		if schema.Value.Properties == nil {
			return objMap, nil
		}

		for fieldName, fieldSchema := range schema.Value.Properties {
			if fieldSchema == nil {
				continue
			}

			currentPath := parentPath
			if currentPath != "" {
				currentPath += "."
			}
			currentPath += fieldName

			var fieldValue interface{}
			var err error
			var fieldType string

			if fieldSchema.Ref != "" {
				parts := strings.Split(strings.TrimPrefix(fieldSchema.Ref, "#/"), "/")
				if len(parts) == 2 && parts[0] == "definitions" {
					if refSchema, ok := swagger.Definitions[parts[1]]; ok {
						fieldValue, err = o.processSchema(swagger, refSchema, currentPath, body)
						if err != nil {
							return nil, fmt.Errorf("failed to process referenced field %s: %w", fieldName, err)
						}
						if refSchema.Value != nil && refSchema.Value.Type != nil && len(refSchema.Value.Type.Slice()) > 0 {
							fieldType = refSchema.Value.Type.Slice()[0]
						}
					}
				}
			} else {
				fieldValue, err = o.processSchema(swagger, fieldSchema, currentPath, body)
				if err != nil {
					return nil, fmt.Errorf("failed to process field %s: %w", fieldName, err)
				}
				if fieldSchema.Value != nil && fieldSchema.Value.Type != nil && len(fieldSchema.Value.Type.Slice()) > 0 {
					fieldType = fieldSchema.Value.Type.Slice()[0]
				}
			}

			objMap[fieldName] = fieldValue

			// 只记录基本类型的字段
			if fieldType != "" && fieldType != "array" && fieldType != "object" {
				body.JSON.Params = append(body.JSON.Params, &HTTPKeyAndType{
					Key:   currentPath,
					Value: "${" + currentPath + "}",
					Type:  fieldType,
				})
			}
		}
		return objMap, nil

	default:
		// 基础类型处理
		if parentPath == "" {
			return "", nil
		}
		// 基本类型直接添加到 Params
		body.JSON.Params = append(body.JSON.Params, &HTTPKeyAndType{
			Key:   parentPath,
			Value: "${" + parentPath + "}",
			Type:  schemaType,
		})
		return "${" + parentPath + "}", nil
	}
}
