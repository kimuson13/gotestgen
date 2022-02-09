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
{{range $tn, $funcName := .TestTargets}}
func Test{{$funcName}}(t *tesitng.T) {
	cases := map[string]struct{
		// write arguments below this

	}{
		// write test cases below this 
		// test case name: {args}

	}

	for testName, tt := range cases {
		tt := tt
		t.Run(testName, func(t *testing.T) {
			// write tests below this
			{{if $.IsParallel}}
			t.Parallel()
			{{end}}
		})
	}
}
{{end}}
`
