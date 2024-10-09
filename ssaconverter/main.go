package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

var LineToVars = make(map[int][]string)

// extractVariablesFromPos extracts all variable names that appear on a specific line.
func InitLineToVars(fset *token.FileSet, node ast.Node) {
	// Inspect the AST of the parsed source code
	ast.Inspect(node, func(n ast.Node) bool {
		// Check if the node is an identifier
		if ident, ok := n.(*ast.Ident); ok {
			pos := fset.Position(ident.Pos())
			LineToVars[pos.Line] = append(LineToVars[pos.Line], ident.Name)
		}
		return true
	})
}

func PrintInstr(instr ssa.Instruction) {
	switch instr.(type) {
	case *ssa.BinOp:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.BinOp).Name(), instr.String())
		position := fset.Position(instr.(*ssa.BinOp).Pos())
		fmt.Println(LineToVars[position.Line])
	case *ssa.Call:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.Call).Name(), instr.String())
	case *ssa.FieldAddr:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.FieldAddr).Name(), instr.String())
	case *ssa.UnOp:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.UnOp).Name(), instr.String())
	case *ssa.Alloc:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.Alloc).Name(), instr.String())
	case *ssa.IndexAddr:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.IndexAddr).Name(), instr.String())
	case *ssa.MakeInterface:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.MakeInterface).Name(), instr.String())
	case *ssa.Slice:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.Slice).Name(), instr.String())
	case *ssa.MakeMap:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.MakeMap).Name(), instr.String())
	case *ssa.Extract:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.Extract).Name(), instr.String())
	case *ssa.MakeChan:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.MakeChan).Name(), instr.String())
	case *ssa.MakeClosure:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.MakeClosure).Name(), instr.String())
	case *ssa.ChangeType:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.ChangeType).Name(), instr.String())
	case *ssa.ChangeInterface:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.ChangeInterface).Name(), instr.String())
	case *ssa.Select:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.Select).Name(), instr.String())
	case *ssa.Lookup:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.Lookup).Name(), instr.String())
	case *ssa.Phi:
		fmt.Printf("%T: %s = %s\n", instr, instr.(*ssa.Phi).Name(), instr.String())
	default:
		fmt.Printf("%T: %s\n", instr, instr.String())
	}
}

var fset *token.FileSet

func main() {
	// 使用flag包来解析命令行参数
	inputFile := flag.String("file", "", "Go source file to parse")
	flag.Parse()

	// 检查是否提供了输入文件
	if *inputFile == "" {
		fmt.Println("Please provide an input file.")
		os.Exit(1)
	}
	// 创建文件集
	// 创建新文件集
	fset = token.NewFileSet()

	// 解析Go源代码
	file, err := parser.ParseFile(fset, *inputFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	InitLineToVars(fset, file) // 建立一个行号到变量名的映射

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
				fmt.Printf("Block: %d\n", block.Index)
				for _, instr := range block.Instrs {
					PrintInstr(instr)
				}
				fmt.Println()
			}
			fmt.Println()
		}
	}
}
