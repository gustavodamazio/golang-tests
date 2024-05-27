package html

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func Titulo(urls ...string) <-chan string {
	c := make(chan string)
	// var wg sync.WaitGroup

	for _, url := range urls {
		// wg.Add(1)
		go func(url string) {
			// defer wg.Done()
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

	// // Fechando o canal quando todas as goroutines terminarem
	// go func() {
	// 	wg.Wait()
	// 	fmt.Println("Fechando o canal", urls)
	// 	close(c)
	// }()

	return c
}
