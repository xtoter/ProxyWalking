package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type parser struct {
	filePath string
	data     []string
}

func NewParser(filepath string) *parser {
	return &parser{filePath: filepath}
}
func (p *parser) GetProxy() *parser {
	file, err := os.Open("data/proxies.txt")
	defer file.Close()
	if err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		p.data = append(p.data, strings.TrimPrefix(fileScanner.Text(), "http://"))
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	return p
}
func (p *parser) Get() []string {
	return p.data
}
