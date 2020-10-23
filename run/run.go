package run

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

//爬虫思路，需要两个channel，一个公共downlaod,每个url对应的解析器对象
//1.首先爬取首页信息，将每个二级子页的url写入channel1
//2.启动gorun爬取每个二级子页，将获取到的图片download,获取到的url写入channel2
//3.最后解析channel2的url，下载对应图片

//首页地址
var index string = "https://www.bizhizu.cn"
//图片存储地址
var savePath string = "C:\\Users\\Administrator\\Desktop\\remake\\go\\import\\src\\go_project\\picture"
//日志存储位置
var logPath string = "C:\\Users\\Administrator\\Desktop\\remake\\go\\import\\src\\go_project\\log"
//并发锁
var wg sync.WaitGroup
//二级子页
type secondPage struct {
	url string
	name string
}
//二级子页的channel
var channel1 = make(chan secondPage,100)

//总入口
func Start() {
	//1.爬取首页信息
	crawIndex()
	//2.遍历channel1,创造对应数量的goroutine来处理
	for i:=0;i<len(channel1);i++ {
		wg.Add(1)
		go secondRead()
	}
	wg.Wait()
}

//1.爬取首页信息
func crawIndex() {
	resp,err := http.Get(index)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic("request fail")
	}
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	doc,err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		panic(err)
	}
	doc.Find(".indextag-list>a").Each(func(i int, selection *goquery.Selection) {
		href,has := selection.Attr("href")
		if has == true {
			//写入channel1中
			channel1 <- secondPage{index+"/"+href,selection.Text()}
		}
	})
	//for i:=0;i<52;i++ {
	//	fmt.Printf("%v\n",<-channel1)
	//}
}

//2.处理首页产生的二级页数据
func secondRead() {
	message := <-channel1
	ok,err := mkmydir(savePath+"\\"+message.name)
	if err != nil {
		MyLog(err.Error()+"\n","secondRead.txt")
	}
	if ok == false {
		MyLog(message.name+"create_dir_err\n","secondRead.txt")
	}
	resp,err := http.Get(message.url)
	if err != nil {
		MyLog(err.Error(),"secondRead.txt")
	}else{
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		MyLog(message.url+"->http_get_code->"+string(resp.StatusCode)+"\n","secondRead.txt")
	}
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		MyLog(err.Error(),"secondRead.txt")
	}
	doc,err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		MyLog(err.Error(),"secondRead.txt")
	}
	doc.Find(".zt_list_left>ul>li>a>img").Each(func(i int, selection *goquery.Selection) {
		alt,_ := selection.Attr("alt")
		src,_ := selection.Attr("src")
		fmt.Printf("start %s--->%s---->%s\n",message.name,alt,src)
		path := savePath+"\\"+message.name+"\\"+alt+".jpg"
		err := DownloadPic(src,path)
		if err != nil {
			MyLog(err.Error(),"secondRead.txt")
		}
		fmt.Printf("end %s--->%s---->%s\n",message.name,alt,src)
	})
	wg.Done()
}
