package common

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpClient struct {
	Host       string
	Port       int32
	AppKey     string
	AppSecret  string
	QueryNum   int
	IsHttps    bool
	HttpClient *http.Client
}

func (hkHttp *HttpClient) Init() {
	//设置tls配置信息
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	hkHttp.HttpClient = &http.Client{Transport: tr}
}

func (hkHttp *HttpClient) HttpGet(path string, headers map[string]string, bodys map[string]string) ([]byte, error) {
	headers = hkHttp.initRequest(headers, path, nil, bodys, false)
	//处理body
	mjson, _ := json.Marshal(bodys)
	mString := string(mjson)

	req, err := http.NewRequest("GET",
		"http://"+hkHttp.Host+":"+NumberToString(hkHttp.Port)+path,
		strings.NewReader(mString))
	if err != nil {
		return nil, err
	}

	//处理header
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		if key != "" {
			req.Header.Set(key, value)
		}
	}
	hkHttp.HttpClient.Timeout = time.Second * 5
	resp, err := hkHttp.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//logs.Info(string(body))
	return body, nil
}

func (hkHttp *HttpClient) HttpPost(path string, headers map[string]string, bodys map[string]string) ([]byte, error) {
	headers = hkHttp.initRequest(headers, path, nil, bodys, true)
	//处理body
	mjson, _ := json.Marshal(bodys)
	mString := string(mjson)

	req, err := http.NewRequest("POST",
		"http://"+hkHttp.Host+":"+NumberToString(hkHttp.Port)+path,
		strings.NewReader(mString))
	if err != nil {
		return nil, err
	}
	req.Method = "POST"

	//处理header
	for key, value := range headers {
		if key != "" {
			req.Header.Set(key, value)
		}
	}
	hkHttp.HttpClient.Timeout = time.Second * 5
	resp, err := hkHttp.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//logs.Info(string(body))
	return body, nil
}

func (hkHttp *HttpClient) HttpPostJson(path string, bodys string) ([]byte, error) {
	header := make(map[string]string)

	accept := "*/*" // "*/*" application/json;
	header["Accept"] = accept

	// ContentType  ;charset=UTF-8
	contentType := "application/json"
	header["Content-Type"] = contentType

	req, err := http.NewRequest("POST",
		path,
		strings.NewReader(bodys))
	if err != nil {
		return nil, err
	}
	req.Method = "POST"

	hkHttp.HttpClient.Timeout = time.Second * 5
	resp, err := hkHttp.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (hkHttp *HttpClient) initRequest(header map[string]string, url string, querys map[string]string, bodys map[string]string, isPost bool) map[string]string {
	// Accept
	accept := "*/*" // "*/*" application/json;
	header["Accept"] = accept

	// ContentType  ;charset=UTF-8
	contentType := "application/json"
	header["Content-Type"] = contentType

	//加上下面的md5，会报验证错误，无法通过验证
	//if isPost==true {
	//	// content-md5，be careful it must be lower case.
	//	mjson, _ := json.Marshal(bodys)
	//	mString := string(mjson)
	//	contentMd5 := computeContentMd5(mString)
	//	header["content-md5"] = contentMd5
	//}
	return header
}
func computeContentMd5(body string) string {
	h := md5.New()
	return base64.StdEncoding.EncodeToString(h.Sum([]byte(body)))
}
