package checker

import (
	"crypto/tls"
	"fmt"
	"github.com/xtoter/ProxyWalking/src/structs"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"
)

type checker struct {
	testAddr string
	proxy    []string
}

var wg sync.WaitGroup

func NewChecker(proxy []string) *checker {
	return &checker{testAddr: "http://ya.ru", proxy: proxy}
}
func (c *checker) CheckActual() []structs.Proxy {
	var result []structs.Proxy
	defaultData := string(getDefaultResponse())
	dataChannel := make(chan structs.Proxy)
	for i, curProxy := range c.proxy {
		wg.Add(1)
		if i%200 == 0 {
			time.Sleep(2 * time.Second)
			fmt.Println(i, "/", len(c.proxy))
		}
		go checkProxy(curProxy, defaultData, dataChannel)
	}
	go func() { // Ждем пока все рутины отработают
		wg.Wait()
		close(dataChannel)
	}()
	for data := range dataChannel {
		result = append(result, data)
	}
	fmt.Println("checked done!")

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time < result[j].Time
	})
	fmt.Println(result)
	return result
}
func getDefaultResponse() []byte {
	requestURL := "https://rdb.altlinux.org/api/version"
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return []byte{}
	}

	// Добавление заголовка accept: application/json
	req.Header.Add("accept", "text/plain")

	// Выполнение запроса
	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return []byte{}
	}
	defer response.Body.Close()

	// Чтение тела ответа в виде среза байтов ([]byte)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа:", err)
		return []byte{}
	}
	return responseBody

}
func checkProxy(proxyAddress string, resultData string, ch chan structs.Proxy) {
	defer wg.Done()
	timeout := 15 * time.Second

	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		fmt.Println("Ошибка при создании SOCKS5-подключения:", err)
		return
	}

	// Создание HTTP-транспорта с настроенным SOCKS5-подключением
	httpTransport := &http.Transport{
		Dial: dialer.Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Пропустить проверку SSL-сертификата
		},
	}

	// Создание клиента HTTP с настроенным транспортом
	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   timeout,
	}

	startTime := time.Now()
	requestURL := "https://rdb.altlinux.org/api/version"
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	// Добавление заголовка accept: application/json
	req.Header.Add("accept", "application/json")

	// Выполнение запроса
	response, err := httpClient.Do(req)
	if err != nil {
		//fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer response.Body.Close()

	// Чтение тела ответа в виде среза байтов ([]byte)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа:", err)
		return
	}
	if string(responseBody) == resultData {
		fmt.Println(proxyAddress, string(responseBody))
		elapsedTime := time.Since(startTime)
		ch <- structs.Proxy{Addr: proxyAddress, Time: elapsedTime, IsBusy: false}
	}

}
