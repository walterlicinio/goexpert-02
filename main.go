package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func fetchAPI(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("URL: %s, Tempo decorrido: %d ms, Resposta: %s", url, int(secs*1000), body)
}
func main() {
	fmt.Println("Digite o CEP:")
	reader := bufio.NewReader(os.Stdin)
	cep, _ := reader.ReadString('\n')
	cep = strings.TrimSpace(cep)

	ch1 := make(chan string)
	ch2 := make(chan string)

	url1 := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	url2 := "http://viacep.com.br/ws/" + cep + "/json/"

	go fetchAPI(url1, ch1)
	go fetchAPI(url2, ch2)

	for {
		select {
		case res := <-ch1:
			fmt.Println(res)
			return
		case res := <-ch2:
			fmt.Println(res)
			return
		case <-time.After(1 * time.Second):
			fmt.Println("Timeout")
			return
		}
	}
}
