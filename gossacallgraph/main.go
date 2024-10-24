package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa/ssautil"
)

// getProjectUsedCall 获取项目使用中的调用方法
func getProjectUsedCall(projectPath string) ([]string, error) {
	projectModule, err := parseProjectModule(projectPath)
	if err != nil {
		return nil, errors.Wrap(err, "parseProjectModule fail")
	}
	log.Printf("projectModule: %+v", projectModule)

	callMap, err := parseProjectCallMap(projectPath)
	if err != nil {
		return nil, errors.Wrap(err, "parseProjectCallMap fail")
	}
	log.Printf("callMap: %+v", callMap)

	srcCall := fmt.Sprintf("%v.main", projectModule)
	isDeleteEdgeFunc := func(caller, callee string) bool {
		// 非本项目调用
		if !strings.Contains(caller, projectModule) || !strings.Contains(callee, projectModule) {
			return true
		}
		// 非初始化调用
		if isInitCall(caller) || isInitCall(callee) {
			return true
		}
		// 非自我调用
		if caller == callee {
			return true
		}
		return false
	}

	// 过滤不需要的边
	for caller, callees := range callMap {
		for callee := range callees {
			if isDeleteEdgeFunc(caller, callee) {
				delete(callees, callee)
			}
		}
		if len(callees) == 0 {
			delete(callMap, caller)
		}
	}

	// 广度搜索图
	for {
		srcCallees := callMap[srcCall]
		srcSize := len(srcCallees)

		for srcCallee := range srcCallees {
			for nextCallee := range callMap[srcCallee] {
				callMap[srcCall][nextCallee] = true
			}
		}

		if srcSize == len(callMap[srcCall]) {
			break
		}
	}

	// 调用源涉及到的所有方法
	var callees []string
	for c := range callMap[srcCall] {
		callees = append(callees, c)
	}
	sort.Strings(callees)
	return callees, nil
}

// parseProjectCallMap 解析项目调用图
func parseProjectCallMap(projectPath string) (map[string]map[string]bool, error) {
	projectModule, err := parseProjectModule(projectPath)
	if err != nil {
		return nil, errors.Wrap(err, "parseProjectModule fail")
	}
	log.Printf("projectModule: %+v", projectModule)

	result, err := analyzeProject(projectPath)
	if err != nil {
		return nil, errors.Wrap(err, "analyzeProject fail")
	}
	log.Printf("analyzeProject: %+v", result)

	// 遍历调用链路
	var callMap = make(map[string]map[string]bool)
	visitFunc := func(edge *callgraph.Edge) error {
		if edge == nil {
			return nil
		}
		// 解析调用者和被调用者
		caller, callee, err := parseCallEdge(edge)
		if err != nil {
			return errors.Wrap(err, "parseCallEdge fail")
		}
		// 记录调用关系
		if callMap[caller] == nil {
			callMap[caller] = make(map[string]bool)
		}
		callMap[caller][callee] = true
		return nil
	}
	err = callgraph.GraphVisitEdges(result.CallGraph, visitFunc)
	if err != nil {
		return nil, errors.Wrap(err, "GraphVisitEdges fail")
	}
	return callMap, nil
}

func parseProjectModule(projectPath string) (string, error) {
	modFilename := filepath.Join(projectPath, "go.mod")
	content, err := ioutil.ReadFile(modFilename)
	if err != nil {
		return "", errors.Wrap(err, "ioutil.ReadFile fail")
	}
	lines := strings.Split(string(content), "\n")
	module := strings.TrimPrefix(lines[0], "module ")
	module = strings.TrimSpace(module)
	return module, nil
}

func analyzeProject(projectPath string) (*pointer.Result, error) {
	// 生成Go Packages
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  projectPath,
	})
	if err != nil {
		return nil, errors.Wrap(err, "packages.Load fail")
	}
	log.Printf("pkgs: %+v", pkgs)

	// 生成ssa 构建编译
	prog, ssaPkgs := ssautil.AllPackages(pkgs, 0)
	prog.Build()
	log.Printf("ssaPkgs: %+v", ssaPkgs)
	// 使用pointer生成调用链路
	return pointer.Analyze(&pointer.Config{
		Mains:          ssaPkgs,
		BuildCallGraph: true,
	})
}

func parseCallEdge(edge *callgraph.Edge) (string, string, error) {
	const callArrow = "-->"
	edgeStr := fmt.Sprintf("%+v", edge)
	strArray := strings.Split(edgeStr, callArrow)
	if len(strArray) != 2 {
		return "", "", fmt.Errorf("invalid format: %v", edgeStr)
	}
	callerNodeStr, calleeNodeStr := strArray[0], strArray[1]
	caller, callee := getCallRoute(callerNodeStr), getCallRoute(calleeNodeStr)
	return caller, callee, nil
}

func getCallRoute(nodeStr string) string {
	nodeStr = strings.TrimSpace(nodeStr)
	if strings.Contains(nodeStr, ":") {
		nodeStr = nodeStr[strings.Index(nodeStr, ":")+1:]
	}
	nodeStr = strings.ReplaceAll(nodeStr, "*", "")
	nodeStr = strings.ReplaceAll(nodeStr, "(", "")
	nodeStr = strings.ReplaceAll(nodeStr, ")", "")
	nodeStr = strings.ReplaceAll(nodeStr, "<", "")
	nodeStr = strings.ReplaceAll(nodeStr, ">", "")
	if strings.Contains(nodeStr, "$") {
		nodeStr = nodeStr[:strings.Index(nodeStr, "$")]
	}
	if strings.Contains(nodeStr, "#") {
		nodeStr = nodeStr[:strings.Index(nodeStr, "#")]
	}
	return strings.TrimSpace(nodeStr)
}

func isInitCall(call string) bool {
	return strings.HasSuffix(call, ".init")
}
