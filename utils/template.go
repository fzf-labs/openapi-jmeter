package utils

import (
	"bytes"
	"html/template"
)

// TemplateExec 执行模板
func TemplateExec(t string, data any) ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(t)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
