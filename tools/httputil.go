package tools

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpDo(method, url string, header map[string]string, body string) (string, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	for k, v := range header {
		request.Header.Add(k, v)
	}

	resp, err := client.Do(request)
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}
