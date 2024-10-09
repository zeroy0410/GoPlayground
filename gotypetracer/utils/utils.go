package utils

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ssa"
)

var LineToVars = make(map[int][]string)

func PrintInstr(instr ssa.Instruction, fset *token.FileSet) {
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

// 打印每个函数的SSA形式
func PrintAllFunc(pkg *ssa.Package, fset *token.FileSet) {
	for _, member := range pkg.Members {
		if fn, ok := member.(*ssa.Function); ok {
			fmt.Printf("Function: %s\n", fn.Name())
			for _, block := range fn.Blocks {
				fmt.Printf("Block: %d\n", block.Index)
				for _, instr := range block.Instrs {
					PrintInstr(instr, fset)
				}
				fmt.Println()
			}
			fmt.Println()
		}
	}
}
