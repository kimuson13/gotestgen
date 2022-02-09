package gotestgen

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"os"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/codegen"
	"github.com/gostaticanalysis/knife"
)

const doc = "gotestgen is test template generate tool"

var (
	flagOutput string
)

func init() {
	Generator.Flags.StringVar(&flagOutput, "o", "", "output file name")
}

var Generator = &codegen.Generator{
	Name: "gotestgen",
	Doc:  doc,
	Run:  run,
}

func run(pass *codegen.Pass) error {
	ifaces := map[string]*knife.Interface{}

	s := pass.Pkg.Scope()
	for _, name := range s.Names() {
		obj := s.Lookup(name)
		if !obj.Exported() {
			continue
		}
		iface, _ := analysisutil.Under(obj.Type()).(*types.Interface)
		if iface != nil {
			ifaces[name] = knife.NewInterface(iface)
		}
	}

	td := &knife.TempalteData{
		Fset:      pass.Fset,
		Files:     pass.Files,
		TypesInfo: pass.TypesInfo,
		Pkg:       pass.Pkg,
	}
	t, err := knife.NewTemplate(td).Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, ifaces); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if flagOutput == "" {
		pass.Print(string(src))
		return nil
	}

	f, err := os.Create(flagOutput)
	if err != nil {
		return err
	}

	fmt.Fprint(f, string(src))

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

var tmpl = `
package {{(pkg).Name}}_test
{{range $tn, $t := .}}
type Mock{{$tn}} struct {
{{- range $n, $f := $t.Methods}}
        {{$n}}Func {{$f.Signature}}
{{- end}}
}
{{range $n, $f := $t.Methods}}
func (m *Mock{{$tn}}) {{$n}}({{range $f.Signature.Params}}
	{{- if (and $f.Signature.Variadic (eq . (last $f.Signature.Params)))}}
        	{{- .Name}} ...{{(slice .Type).Elem}},
	{{- else}}
        	{{- .Name}} {{.Type}},
	{{- end}}
{{- end}}) ({{range $f.Signature.Results}}
        {{- .Name}} {{.Type}},
{{- end}}) {
        {{if $f.Signature.Results}}return {{end}}m.{{$n}}Func({{range $f.Signature.Params}}
		{{- if (and $f.Signature.Variadic (eq . (last $f.Signature.Params)))}}
        		{{- .Name}}...,
		{{- else}}
        		{{- .Name}},
		{{- end}}
        {{- end}})
}
{{end}}
{{end}}
`
