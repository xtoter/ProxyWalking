package main

import (
	"fmt"
	"github.com/xtoter/ProxyWalking/src/checker"
	"github.com/xtoter/ProxyWalking/src/handler"
	"github.com/xtoter/ProxyWalking/src/parser"
)

func main() {
	proxy := parser.NewParser("data/proxies.txt").GetProxy().Get()
	actualProcy := checker.NewChecker(proxy).CheckActual()
	fmt.Println(len(actualProcy))
	handler.NewHandler(actualProcy).Run()
}
