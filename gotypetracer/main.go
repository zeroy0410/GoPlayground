package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"gotypetracer/utils"
	"log"
	"os"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

var fset *token.FileSet

func main() {
	inputFile := flag.String("file", "", "Go source file to parse")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Please provide an input file.")
		os.Exit(1)
	}

	fset = token.NewFileSet()

	// 解析Go源代码
	file, err := parser.ParseFile(fset, *inputFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	utils.InitLineToVars(fset, file) // 建立一个行号到变量名的映射

	// 设置类型检查器配置
	conf := loader.Config{Fset: fset}
	conf.CreateFromFiles("example", file)

	// 加载程序
	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	// 创建SSA程序
	ssaProg := ssautil.CreateProgram(prog, ssa.SanityCheckFunctions)
	ssaProg.Build()

	// 获取包并生成SSA
	pkg := ssaProg.Package(prog.Created[0].Pkg)
	pkg.Build()

	utils.PrintAllFunc(pkg, fset)
}
