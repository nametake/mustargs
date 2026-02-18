package main

import (
	"github.com/nametake/mustargs"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(mustargs.Analyzer) }
