package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

// Google I/O 2012 - Go Concurrency Patterns

// <-chan - canal somente-leitura
func titulo(urls ...string) <-chan string {
	c := make(chan string)
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				c <- fmt.Sprintf("Erro ao acessar %s: %v", url, err)
				return
			}
			defer resp.Body.Close()
			html, err := io.ReadAll(resp.Body)
			if err != nil {
				c <- fmt.Sprintf("Erro ao ler %s: %v", url, err)
				return
			}

			r, _ := regexp.Compile("<title>(.*?)<\\/title>")
			matches := r.FindStringSubmatch(string(html))
			if len(matches) > 1 {
				c <- matches[1]
			} else {
				c <- fmt.Sprintf("Título não encontrado em %s", url)
			}
		}(url)
	}

	// Fechando o canal quando todas as goroutines terminarem
	go func() {
		wg.Wait()
		fmt.Println("Fechando o canal", urls)
		close(c)
	}()

	return c
}

func main() {
	t1 := titulo("https://www.cod3r.com.br", "https://www.google.com")
	t2 := titulo("https://www.amazon.com", "https://www.youtube.com")

	fmt.Println("Primeiros:", <-t1, "|", <-t2)
	fmt.Println("Segundos:", <-t1, "|", <-t2)
	// Ler do canal até que ele seja fechado
	// for titulo := range t1 {
	// 	fmt.Println("Titulo:", titulo)
	// }

	// for titulo := range t2 {
	// 	fmt.Println("Titulo:", titulo)
	// }
}
