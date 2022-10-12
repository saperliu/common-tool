package dd

import (
	"crypto/tls"
	"github.com/saperliu/common-tool/logger"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type DDHookMessage struct {
	HookUrl string
}

func (opd *DDHookMessage) SendActionMsg(msg string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("-----send ding ding msg  error panic :  %v ", err)
		}
	}()
	headers := make(map[string]string)

	headers["Content-Type"] = "application/json ;charset=utf-8"
	dataTime := time.Now().Format("2006-01-02 15:04:05")
	content := "{\"msgtype\":\"text\",\"text\":{\"content\": \"系统监控: " + msg + "[" + dataTime + "]\"}}"
	//产品部群组 "https://oapi.dingtalk.com/robot/send?access_token=dedb7bb4d9d0775d29e5e95c6066"
	req, err := http.NewRequest("POST", opd.HookUrl,
		strings.NewReader(content))

	for key, value := range headers {
		if req != nil && key != "" {
			req.Header.Set(key, value)
		}
	}

	//设置tls配置信息
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	httpClient := &http.Client{Transport: tr}
	httpClient.Timeout = time.Second * 5

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("request url  error : %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("request body read all   error : %v", err)
	}
	result := string(body)
	logger.Info("response is %#v\n", result)

}
