package run

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

//创建文件夹
func mkmydir(path string) (bool,error) {
	err := os.MkdirAll(path,os.ModePerm)
	if err != nil {
		return false,err
	}
	return true,nil
}

//日志文件
func MyLog(message string,filename string) error {
	file,err := os.OpenFile(logPath+"\\"+filename,os.O_WRONLY|os.O_APPEND|os.O_CREATE,os.ModePerm)
	if err != nil {
		return err
	}
	_,err = file.WriteString(message+"\n")
	if err != nil {
		return err
	}
	return nil
}

//下载图片 https://uploadfile.bizhizu.cn/up/46/3c/73/463c7345447f0f14eb41931269ee46f1.jpg
func DownloadPic(url string,path string) error {
	resp,err := http.Get(url)
	if err != nil {
		return err
	}
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	out,err := os.Create(path)
	if err != nil {
		return err
	}
	io.Copy(out,bytes.NewReader(body))
	return nil
}
