package main

import (
	"bytes"
	"strings"
	"text/template"
)

var httpTemplate = `
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
type {{.ServiceType}}ProtoSpec struct {
	Method string
	HttpPath string
	HttpAction string
}

type {{.ServiceType}}APISpec struct {
	Service string
	Methods []*{{.ServiceType}}ProtoSpec 
}

func New{{.ServiceType}}APISpec() interface{} {
	return &{{.ServiceType}}APISpec{
		Service: "{{$svrName}}",
		Methods: []*{{.ServiceType}}ProtoSpec{
{{- range .Methods}}
			{
				"{{.Name}}",
				"{{.Path}}",
				"{{.Method}}",
			},
{{- end}}
		},
	}
}
`

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/helloworld/helloworld.proto
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc
}

type methodDesc struct {
	// method
	Name    string
	Num     int
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	HasVars      bool
	HasBody      bool
	Body         string
	ResponseBody string
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
