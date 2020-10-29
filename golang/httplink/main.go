package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//func CurlGet(url string, timeout time.Duration) (result []byte, err error) {
//	cli := &http.Client{}
//
//	req, err := http.NewRequest(http.MethodGet, url, nil)
//	if err != nil {
//		err = errors.New("err req")
//		return
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), timeout)
//	defer cancel()
//
//	resp, err := cli.Do(req)
//	if err != nil {
//
//	}
//	defer resp.Body.Close()
//
//	result, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//
//	}
//
//	return
//}
//
//func CurlFormPost(uri, token string, params map[string]interface{}, timeout time.Duration) (result []byte, err error) {
//	cli := &http.Client{}
//
//	values := url.Values{}
//	for k, v := range params {
//		if v != nil {
//			values.Set(k, cast.ToString(v))
//		}
//	}
//
//	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(values.Encode()))
//
//	if err != nil {
//
//	}
//
//	req.Header.Set("ACCESS-TOKEN", token)
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//}

func main() {

	cli := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
	if err != nil {
		return
	}
	resp, err := cli.Do(req)
	if err != nil {
		return
	}
	// 必须关闭, 如果我们没有写关闭resp.Body打开的句柄,就会导致句柄泄露
	// defer resp.Body.Close() //
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(data))
	return
}
