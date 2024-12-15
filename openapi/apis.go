package openapi

type API struct {
	HTTPName    string            `json:"httpName"`    // http请求名称
	HTTPDomain  string            `json:"httpDomain"`  // http请求域名
	HTTPPath    string            `json:"httpPath"`    // http请求路径
	HTTPMethod  string            `json:"httpMethod"`  // http请求方法
	HTTPHeaders []*HTTPKeyAndType `json:"httpHeaders"` // http请求头
	HTTPParams  HTTPParams        `json:"httpParams"`  // http请求参数
	HTTPBody    HTTPBody          `json:"httpBody"`    // http请求体
}
type HTTPKeyAndType struct {
	Key   string `json:"key"`   // 键
	Value string `json:"value"` // 值
	Type  string `json:"type"`  // 类型
}

type HTTPParams struct {
	Query []*HTTPKeyAndType `json:"query"` // http 请求查询参数
	Path  []*HTTPKeyAndType `json:"path"`  // http 请求路径参数
}

type HTTPBody struct {
	ContentType string            `json:"contentType"` // http 请求内容类型 form-data,application/x-www-form-urlencoded,application/json
	FormData    []*HTTPKeyAndType `json:"formData"`    // http 请求表单数据
	JSON        *HTTPBodyJSON     `json:"json"`        // http 请求 json 数据
}

type HTTPBodyJSON struct {
	JSONStr string            `json:"jsonStr"` // http 请求 json 数据
	Params  []*HTTPKeyAndType `json:"params"`  // http 请求 json 参数
}
