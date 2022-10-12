package mail

import (
	"common-tool/logger"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"os/exec"
)

func SendMail(title string, content string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("-----SendMail error panic :  %v ", err)
		}
	}()

	sysInfo := getSystemInfo()

	d := gomail.NewDialer("smtp.163.com", 25, "test@163.com", "123456")
	s, err := d.Dial()
	logger.Info("---------send email s------- %v", s)
	if err != nil {
		//panic(err)
		logger.Error("---------send email error------- %v", err)
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "test@163.com")
	m.SetHeader("To", "saperliu@163.com")
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content+sysInfo)
	if err := gomail.Send(s, m); err != nil {
		logger.Error("---------send email error------- %v", err)
	}
	logger.Info("---------send   email ------- %v  %v ", title, content)
}

func getSystemInfo() string {
	cmd := exec.Command("netstat", "-ano")
	stdout, err := cmd.StdoutPipe()
	if err != nil { //获取输出对象，可以从该对象中读取输出结果
		logger.Error("-----backRun StdoutPipe :  %v ", err)
	}
	defer stdout.Close()                // 保证关闭输出流
	if err := cmd.Start(); err != nil { // 运行命令
		logger.Error("-----backRun cmd.Start() :  %v ", err)
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil { // 读取输出结果
		logger.Error("-----opBytes  err :  %v ", err)
	} else {
		logger.Info("-----backRun string(opBytes) :  %v ", string(opBytes))
	}
	return string(opBytes)
}
