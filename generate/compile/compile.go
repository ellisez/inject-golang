package compile

import (
	"bytes"
	"text/template"
)

func executeTemplate(tmpl *template.Template, data any) (string, error) {
	buffer := &bytes.Buffer{}
	err := tmpl.Execute(buffer, data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
