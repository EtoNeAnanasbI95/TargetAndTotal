package main

import (
	"fmt"
	TargetAndTotal "github.com/EtoNeAnanasbI95/TargetAndTotal/internal/app"
	"github.com/EtoNeAnanasbI95/TargetAndTotal/internal/ollama"
	"github.com/fatih/color"
	"log"
	"sync"
)

var (
	title     = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	highlight = color.New(color.FgHiCyan, color.Bold).SprintFunc()
	loading   = color.New(color.FgYellow, color.Italic).SprintFunc()
	resColor  = color.New(color.FgHiGreen, color.Bold).SprintFunc()
	alert     = color.New(color.FgHiRed, color.Bold).SprintFunc()
)

func main() {
	errChan := make(chan error, 1)
	resChan := make(chan string)
	wg := sync.WaitGroup{}

	defer close(resChan)
	defer close(errChan)

	lama := ollama.NewOllama("llama3.2", "http://localhost:11434/api/generate", false, "application/json")
	app := TargetAndTotal.NewApp(lama)
	app.Run(&wg, errChan, resChan)
	var result string

	select {
	case err := <-errChan:
		fmt.Println()
		log.Fatal(err.Error())
	case res := <-resChan:
		result = res
	}
	wg.Wait()
	fmt.Println()
	fmt.Println(resColor(result))
}
