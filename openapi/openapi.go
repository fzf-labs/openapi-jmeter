package openapi

import (
	"fmt"

	"github.com/fzf-labs/openapi-jmeter/config"
)

type OpenAPI struct {
	conf *config.Config
}

func NewOpenAPI(conf *config.Config) *OpenAPI {
	return &OpenAPI{
		conf: conf,
	}
}

func (o *OpenAPI) Run() ([]*API, error) {
	switch o.conf.Jmeter.OpenapiVersion {
	case "2.0":
		return NewOpenAPI2(o.conf).Run()
	case "3.0":
		return NewOpenAPI3(o.conf).Run()
	default:
		return nil, fmt.Errorf("invalid OpenAPI version: %s", o.conf.Jmeter.OpenapiVersion)
	}
}
