package main

import (
	"ProxyWalking/src/checker"
	"ProxyWalking/src/handler"
	"ProxyWalking/src/parser"
	"fmt"
)

func main() {
	proxy := parser.NewParser("data/proxies.txt").GetProxy().Get()
	actualProcy := checker.NewChecker(proxy).CheckActual()
	fmt.Println(len(actualProcy))
	handler.NewHandler(actualProcy).Run()
}
