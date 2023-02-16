package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type JsonWrapper struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Msg     string `json:"msg"`
}

type Request struct {
	Url    string
	Token  string
	Params url.Values
}

// 发送post请求，并解析key的数据
func (request Request) PostParse(v any) (string, error) {
	message := "网络请求失败"

	if request.Url == "" {
		message = "url 不能为空"
		return message, errors.New(message)
	}

	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, request.Url, strings.NewReader(request.Params.Encode()))
	if err != nil {
		return message, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req) //http.Post(request.Url, "application/x-www-form-urlencoded", strings.NewReader(request.Params.Encode()))
	if err != nil {
		return message, err
	}

	defer resp.Body.Close()
	//Read the response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return message, err
	}

	if resp.StatusCode == http.StatusNotFound {
		message = "404"
		err = errors.New(message)
		return message, err
	} else if resp.StatusCode == http.StatusUnauthorized {
		message = "未登录"
		err = errors.New(message)
		return message, err
	}
	fmt.Printf("string(body): %v\n", string(body))
	err = json.Unmarshal(body, v)
	if err != nil {
		message = "数据解析失败"
		return message, err
	}

	return "", nil
}
