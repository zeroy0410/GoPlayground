package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	// 创建文件集
	// 创建新文件集
	fset := token.NewFileSet()

	// 解析Go源代码
	file, err := parser.ParseFile(fset, "your_file.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

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

	// 打印每个函数的SSA形式
	for _, member := range pkg.Members {
		if fn, ok := member.(*ssa.Function); ok {
			fmt.Printf("Function: %s\n", fn.Name())
			for _, block := range fn.Blocks {
				for _, instr := range block.Instrs {
					fmt.Printf("%T: %s\n", instr, instr.String())
				}
			}
			fmt.Println()
		}
	}
}
