package handler

import (
	"crypto/tls"
	"fmt"
	"github.com/xtoter/ProxyWalking/src/structs"
	"golang.org/x/net/proxy"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	proxy []structs.Proxy
}

func NewHandler(proxy []structs.Proxy) *Handler {
	return &Handler{proxy: proxy}
}

func (h *Handler) Run() {
	http.HandleFunc("/", h.handleRequest)
	fmt.Println("HTTP server started on :8080")
	http.ListenAndServe(":8080", nil)
}
func (h *Handler) getProxy() int {
	for i, cur := range h.proxy {
		if !cur.IsBusy {
			fmt.Println(i)
			return i
		}
	}
	return -1
}
func (h *Handler) handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error

	var response *http.Response
	cond := true
	for cond {
		err, response = h.goToOriginal(w, r)
		if err == nil {
			fmt.Println("sucsess")
			cond = false
		} else {
			fmt.Println("error,repeat ")
		}
	}

	defer response.Body.Close()
	// Копируем заголовки ответа от удаленного сервера к клиенту
	copyHeaders(w.Header(), response.Header)

	// Отправляем статус и тело ответа клиенту
	w.WriteHeader(response.StatusCode)

	io.Copy(w, response.Body)
	//

}
func (h *Handler) goToOriginal(w http.ResponseWriter, r *http.Request) (error, *http.Response) {
	// Создаем клиента с настройками прокси
	proxyInd := h.getProxy()
	if proxyInd == -1 {
		return fmt.Errorf("%s", "Error get proxy"), nil
	}
	h.proxy[proxyInd].IsBusy = true
	fmt.Println(h.proxy[proxyInd].Addr)
	dialer, err := proxy.SOCKS5("tcp", h.proxy[proxyInd].Addr, nil, proxy.Direct)
	if err != nil {

		return fmt.Errorf("%s", "Ошибка при создании SOCKS5-подключения:"), nil
	}

	// Создание HTTP-транспорта с настроенным SOCKS5-подключением
	httpTransport := &http.Transport{
		Dial: dialer.Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Пропустить проверку SSL-сертификата
		},
	}

	// Создание клиента HTTP с настроенным транспортом
	client := &http.Client{
		Transport: httpTransport,
		Timeout:   time.Second * 25,
	}

	//client = &http.Client{}
	//
	// Создаем новый запрос для перенаправления
	requestURL := "https://rdb.altlinux.org" + r.URL.Path
	request, err := http.NewRequest(r.Method, requestURL, r.Body)
	if err != nil {
		return fmt.Errorf("%s", "Error creating request:"), nil
	}
	// Аргументы запроса
	request.URL.RawQuery = r.URL.RawQuery
	// Копируем заголовки запроса клиента к новому запросу
	copyHeaders(request.Header, r.Header)

	// Выполняем запрос на удаленный сервер
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("%s", "Error forwarding request:"), nil
	}
	go func() {
		time.Sleep(time.Second * 1)

		h.proxy[proxyInd].IsBusy = false
	}()
	return nil, response
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
