package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type ApiCepResponse struct {
	Status   int    `json:"status"`
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func fetchApiCep(cep string, ch chan<- string) {
	cep = cep[:5] + "-" + cep[5:]
	start := time.Now()
	url := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()
	var apiCepResponse ApiCepResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	err = json.Unmarshal(body, &apiCepResponse)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("URL: %s, Tempo decorrido: %d ms, Resposta: %+v", url, int(secs*1000), apiCepResponse)
}

func fetchViaCep(cep string, ch chan<- string) {
	start := time.Now()
	url := "http://viacep.com.br/ws/" + cep + "/json/"
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()
	var viaCepResponse ViaCepResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	err = json.Unmarshal(body, &viaCepResponse)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("URL: %s, Tempo decorrido: %d ms, Resposta: %s", url, int(secs*1000), viaCepResponse)
}

func main() {
	fmt.Println("Digite o CEP:")
	reader := bufio.NewReader(os.Stdin)
	cep, _ := reader.ReadString('\n')
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	cep = reg.ReplaceAllString(cep, "")
	if len(cep) != 8 {
		fmt.Println("CEP precisa ter 8 dígitos.")
		return
	}

	ch1 := make(chan string)
	ch2 := make(chan string)

	go fetchApiCep(cep, ch1)
	go fetchViaCep(cep, ch2)

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
