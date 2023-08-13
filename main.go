package main

import (
	"fmt"
	"github.com/xtoter/ProxyWalking/src/checker"
	"github.com/xtoter/ProxyWalking/src/handler"
	"github.com/xtoter/ProxyWalking/src/parser"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Укажите файл с proxy")
		return
	}
	proxy := parser.NewParser(os.Args[1]).GetProxy().Get()
	actualProcy := checker.NewChecker(proxy).CheckActual()
	fmt.Println("Актуальных проксей ", len(actualProcy))
	handler.NewHandler(actualProcy).Run()
}
