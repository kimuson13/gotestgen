package main

import (
	"github.com/gostaticanalysis/codegen/singlegenerator"
	"github.com/kimuson13/gotestgen"
)

func main() {
	singlegenerator.Main(gotestgen.Generator)
}
