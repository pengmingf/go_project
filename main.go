package main

import (
	"fmt"
	"go_project/run"
)

func main() {
	fmt.Println("start crawler")
	run.Start()
	//err := run.MyLog("test\n","test.txt")
	//if err != nil {
	//	panic(err)
	//}
	//run.DownloadPic("https://uploadfile.bizhizu.cn/up/46/3c/73/463c7345447f0f14eb41931269ee46f1.jpg","C:\\Users\\Administrator\\Desktop\\remake\\go\\import\\src\\go_project\\picture\\1.jpg")
	fmt.Println("the crawler has end")
}
