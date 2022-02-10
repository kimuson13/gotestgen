package gotestgen

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"os"
	"strings"
	"unicode"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/codegen"
	"github.com/gostaticanalysis/knife"
)

const doc = "gotestgen is test template generate tool"

var (
	flagIsParallel        bool
	flagGenerateFilePaths string
)

var flagDesc string = `
["package name":"filepath" "other package":"filepath"]
filepath accept only directory
please see github.com/kimuson13/gotestgen to know more info.
`

func init() {
	Generator.Flags.BoolVar(&flagIsParallel, "p", false, "whether t.Parallel() or not")
	Generator.Flags.StringVar(&flagGenerateFilePaths, "g", "", flagDesc)
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
	testTargets := make(map[types.Object]string)

	for key, val := range pass.TypesInfo.Defs {
		switch val.(type) {
		case *types.Func:
			if r := rune(key.Name[0]); unicode.IsUpper(r) {
				testTargets[val] = key.Name
			}
		}
	}

	var fileName string
	for _, v := range pass.Files {
		pkgName := v.Name.Name
		if pkgName == "main" {
			return nil
		} else if strings.HasSuffix(pkgName, "_test") {
			return nil
		}

		fileName = pkgName
	}

	s := pass.Pkg.Scope()
	for _, name := range s.Names() {
		obj := s.Lookup(name)
		if !obj.Exported() {
			continue
		}
		iface, ok := analysisutil.Under(obj.Type()).(*types.Interface)
		if !ok {
			continue
		}
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

	f, err := os.OpenFile(fmt.Sprintf("%s_test.go", fileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprint(f, string(src)); err != nil {
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
