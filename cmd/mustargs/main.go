package main

import (
	"github.com/nametake/mustargs"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(mustargs.Analyzer) }
