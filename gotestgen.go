package gotestgen

import (
	"bytes"
	"go/format"
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/codegen"
	"github.com/gostaticanalysis/knife"
)

const doc = "gotestgen is test template generate tool"

var (
	flagIsParallel bool
)

func init() {
	Generator.Flags.BoolVar(&flagIsParallel, "p", false, "whether t.Parallel() or not")
}

var Generator = &codegen.Generator{
	Name: "gotestgen",
	Doc:  doc,
	Run:  run,
}

type ExecuteData struct {
	TestTargets map[types.Object]string
	IsParallel  bool
}

func run(pass *codegen.Pass) error {
	testTargets := map[types.Object]string{}

	for key, val := range pass.TypesInfo.Defs {
		switch val.(type) {
		case *types.Func:
			testTargets[val] = key.Name
		}
	}

	s := pass.Pkg.Scope()
	for _, name := range s.Names() {
		obj := s.Lookup(name)
		if !obj.Exported() {
			continue
		}
		iface, _ := analysisutil.Under(obj.Type()).(*types.Interface)
		for i := 0; i < iface.NumMethods(); i++ {
			if _, ok := testTargets[iface.Method(i)]; ok {
				delete(testTargets, iface.Method(i))
			}
		}
	}

	ed := ExecuteData{TestTargets: testTargets, IsParallel: flagIsParallel}

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
	if err := t.Execute(&buf, ed); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if _, err := pass.Print(string(src)); err != nil {
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
