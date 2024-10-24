package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const hello = `
package main

import (
	"fmt"
	"time"
)

// 定义一个结构体
type Person struct {
	Name string
	Age  int
}

type Student struct {
	Name string
	Age  int
}

// 定义一个接口
type Worker interface {
	Work()
}

// Person 实现了 Worker 接口
func (p Person) Work() {
	fmt.Printf("%s is working...\n", p.Name)
}

// Student 实现了 Worker 接口
func (s Student) Work() {
	fmt.Printf("%s is working...\n", s.Name)
}

// 使用 map 存储工人的工资
var salaries = map[Worker]int{
	Person{Name: "Alice", Age: 30}: 50000,
	Person{Name: "Bob", Age: 25}:   40000,
}

// 使用 select 来处理多个 channel
func handleCommunications(c1, c2 <-chan string) {
	for {
		select {
		case msg1 := <-c1:
			fmt.Println("Received from channel 1:", msg1)
		case msg2 := <-c2:
			fmt.Println("Received from channel 2:", msg2)
		case <-time.After(5 * time.Second):
			fmt.Println("Timed out waiting for messages.")
			return
		}
	}
}

func Run(p Worker) {
	p.Work()
}

func main() {
	// 创建两个 channel
	ch1 := make(chan string)
	ch2 := make(chan string)
	x := "I am X!"

	// 启动一个 goroutine 来发送消息到 channel 1
	go func() {
		fmt.Println(x)
		ch1 <- "Hello from channel 1"
	}()

	// 启动一个 goroutine 来发送消息到 channel 2
	go func() {
		ch2 <- "Hello from channel 2"
	}()

	// 启动一个 goroutine 来处理通信
	go handleCommunications(ch1, ch2)

	// 创建一个 Worker 实例
	worker := Person{Name: "Charlie", Age: 35}
	student := Student{Name: "Student", Age: 18}
	Run(worker)
	Run(student)

	// 打印工资
	if salary, ok := salaries[worker]; ok {
		fmt.Printf("%s's salary is $%d\n", worker.Name, salary)
	} else {
		fmt.Printf("%s's salary is not recorded\n", worker.Name)
	}

	// 为了等待所有 goroutine 完成，我们等待一段时间
	time.Sleep(10 * time.Second)
}

`

func main() {
	// Parse the source files.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", hello, parser.ParseComments)
	if err != nil {
		fmt.Print(err) // parse error
		return
	}
	files := []*ast.File{f}

	// Create the type-checker's package.
	pkg := types.NewPackage("hello", "")

	// Type-check the package, load dependencies.
	// Create and build the SSA program.
	hello, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()}, fset, pkg, files, ssa.SanityCheckFunctions)
	if err != nil {
		fmt.Print(err) // type error in some package
		return
	}

	// Print out the package.
	hello.WriteTo(os.Stdout)

	// Print out the package-level functions.
	hello.Func("init").WriteTo(os.Stdout)
	hello.Func("main").WriteTo(os.Stdout)
}
