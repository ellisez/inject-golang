package compile

import (
	"fmt"
	"generate/model"
	"text/template"
)

func init() {
	var err error

	singletonTmpl, err = template.New("singletonStruct").Parse(singletonTmplText)
	if err != nil {
		panic(err)
	}

	constructorTmpl, err = template.New("localStruct").Parse(constructorTmplText)
	if err != nil {
		panic(err)
	}

	multipleTmpl, err = template.New("multipleStruct").Parse(multipleTmplText)
	if err != nil {
		panic(err)
	}
}

var singletonTmpl *template.Template
var singletonTmplText = `
package {{.Package}}
import (
	. "inject"
	"{{.Dirname}}"
)
func init() {
	Set[*{{.Package}}.{{.Name}}]("{{.Instance}}", &{{.Package}}.{{.Name}}{})
}
`

// {{- range $index, $element := .NormalFields -}}
// {{- $element.Name}} {{$element.Type -}}
// {{if gt $index 0}}, {{end}}
// {{- end -}}
var constructorTmpl *template.Template
var constructorTmplText = `
package {{.Package}}
import (
	. "inject"
)
func {{.Constructor}}() *{{.Name}} {
	return Get[*{{.Name}}]("{{.Instance}}")
}`

var multipleTmpl *template.Template
var multipleTmplText = `
package {{.Package}}
import (
	. "inject"
)
func {{.Constructor}}(
	{{- range $index, $element := .NormalFields -}}
		{{- $element.Name}} {{$element.Type -}}
		{{if gt $index 0}}, {{end}}
	{{- end -}}
) *{{.Name}} {
	{{.Instance}} := &{{.Name}}{
		{{- range $element := .NormalFields -}}
			{{- $element.Name}}: {{$element.Name -}},
		{{- end -}}
		{{- range $element := .InjectFields -}}
			{{- $element.Name}}: Get[*{{$element.Name}}]("{{$element.Instance}}"),
		{{- end -}}
	}
	return {{.Instance}}
}`

func SingletonOfStruct(structAnnotate *model.StructInfo) (string, error) {
	genCode, err := executeTemplate(singletonTmpl, structAnnotate)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	}
	return genCode, nil
}

func LocalOfStruct(structAnnotate *model.StructInfo) (string, error) {
	genCode, err := executeTemplate(constructorTmpl, structAnnotate)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	}
	return genCode, nil
}
func MultipleOfStruct(structAnnotate *model.StructInfo) (string, error) {
	genCode, err := executeTemplate(multipleTmpl, structAnnotate)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	}
	return genCode, nil
}
