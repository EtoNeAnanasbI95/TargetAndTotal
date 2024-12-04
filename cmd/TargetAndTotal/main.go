package main

import (
	"TargetAndTotal/internal/ollama"
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	title     = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	highlight = color.New(color.FgHiCyan, color.Bold).SprintFunc()
	loading   = color.New(color.FgYellow, color.Italic).SprintFunc()
	resColor  = color.New(color.FgHiGreen, color.Bold).SprintFunc()
	alert     = color.New(color.FgHiRed, color.Bold).SprintFunc()
)

func main() {
	fmt.Println(title("Добро пожаловать в генератор текста!"))
	fmt.Println(alert("Для выхода из программа нажмите ctrl + C"))
	fmt.Println()
	fmt.Println("Введите текст из технического задания. Для завершения ввода введите пустую строку или нажмите Ctrl+D (Linux/Mac) / Ctrl+Z (Windows):")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder
	for {
		fmt.Print(highlight("* "))
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		builder.WriteString(line)
	}

	input := builder.String()
	message := strings.Replace(input, "\n", " ", -1)

	errChan := make(chan error, 1)
	resChan := make(chan string)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	defer close(resChan)
	defer close(errChan)

	lama := ollama.NewOllama("llama3.2", "http://localhost:11434/api/generate", false, "application/json")

	wg.Add(1)
	go func(func()) {
		prompt := fmt.Sprintf("Сформулируй краткую цель работы без перечислений исходя из вот этого текста \\\"%v\\\", твой ответ должен начинаться с \\\"Цель: ...\\\" и быть длиной 40 слов", message)
		target, err := lama.Generate(prompt)
		if err != nil {
			errChan <- err
			return
		}
		prompt = fmt.Sprintf("Сформулируй краткий вывод работы без перечислений исходя из вот этой цели \\\"%v\\\", твой ответ должен начинаться с \\\"Вывод: ...\\\" и быть длиной 40 слов", target)
		total, err := lama.Generate(prompt)
		if err != nil {
			errChan <- err
			return
		}
		cancel()
		resChan <- fmt.Sprintf("%v \n\n...\n\n %v", target, total)
		wg.Done()
	}(cancel)

	fmt.Println()
	fmt.Println()

	wg.Add(1)
	go func(ctx context.Context) {
		frames := []string{"-", "\\", "|", "/"}
		i := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Print("\r")
				wg.Done()
				return
			default:
				fmt.Print(loading(fmt.Sprintf("\rЗагрузка... %s", frames[i%len(frames)])))
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}(ctx)

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
