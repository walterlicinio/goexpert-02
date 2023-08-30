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
	ch <- fmt.Sprintf("%.2f ms decorrido com resposta da API %s\nResposta: %s", secs*1000, url, body)
}
func main() {
	fmt.Println("Digite o CEP:")
	reader := bufio.NewReader(os.Stdin)
	cep, _ := reader.ReadString('\n')
	cep = strings.TrimSpace(cep)

	ch := make(chan string)

	url1 := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	url2 := "http://viacep.com.br/ws/" + cep + "/json/"

	go fetchAPI(url1, ch)
	go fetchAPI(url2, ch)

	timeout := time.After(1 * time.Second)
	select {
	case res := <-ch:
		fmt.Println(res)
	case <-timeout:
		fmt.Println("Timeout")
	}
}
