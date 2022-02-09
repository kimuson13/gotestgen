package gotestgen_test

import (
	"flag"
	"os"
	"testing"

	"github.com/gostaticanalysis/codegen/codegentest"
	"github.com/kimuson13/gotestgen"
)

var flagUpdate bool

func TestMain(m *testing.M) {
	flag.BoolVar(&flagUpdate, "update", false, "update the golden files")
	flag.Parse()
	os.Exit(m.Run())
}

func TestGenerator(t *testing.T) {
	rs := codegentest.Run(t, codegentest.TestData(), gotestgen.Generator, "a")
	codegentest.Golden(t, rs, flagUpdate)
}
