package gotestgen_test

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/gostaticanalysis/codegen/codegentest"
	"github.com/kimuson13/gotestgen"
)

var flagUpdate bool

func TestMain(m *testing.M) {
	flag.BoolVar(&flagUpdate, "update", false, "update the golden files")
	flag.Parse()
	os.Exit(removeAndExit(m.Run()))
}

func removeAndExit(code int) int {
	if err := os.Remove("a_test.go"); err != nil {
		log.Println("failed to remove file: ", err)
	}

	return code
}

func TestGenerator(t *testing.T) {
	rs := codegentest.Run(t, codegentest.TestData(), gotestgen.Generator, "a")
	codegentest.Golden(t, rs, flagUpdate)
}
