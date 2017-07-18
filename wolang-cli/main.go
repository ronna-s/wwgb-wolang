package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/julianbachiller/wwgb-wolang/wolang"
)

func startClient() (chan string, chan interface{}, chan error) {

	inputCh := make(chan string)
	resultCh := make(chan interface{})
	errCh := make(chan error)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("wolang> ")
			if input, readErr := reader.ReadString('\n'); readErr != nil {
				panic(readErr.Error())
			} else {
				inputCh <- input
			}
			select {
			case err := <-errCh:
				fmt.Println(err.Error())
			case result := <-resultCh:
				fmt.Println(result)
			}
		}
	}()
	return inputCh, resultCh, errCh;
}

func startParsing(inputCh chan string, errCh chan error) chan wolang.DataType {
	parsedCh := make(chan wolang.DataType)
	go func() {
		var (
			err error
			parsed wolang.DataType
		)
		for input := range inputCh {
			_, parsed, err = wolang.Parse(strings.TrimSuffix(strings.TrimSuffix(input, "\n"), "\r"))
			if err != nil {
				fmt.Println(err)
				errCh <- err
			}
			parsedCh <- parsed
		}
	}()
	return parsedCh
}

func startEvaluating(parsedCh chan wolang.DataType, resultCh chan interface{}, errCh chan error) {
	go func() {
		for parsed := range parsedCh {
			if result, err := wolang.Eval(parsed); err != nil {
				errCh <- err
			} else {
				resultCh <- result
			}
		}
	}()
}

func main() {

	readInputCh, resultCh, errCh := startClient()
	readParseDataCh := startParsing(readInputCh, errCh)
	startEvaluating(readParseDataCh, resultCh, errCh)

	for true {
	}
}
